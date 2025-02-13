package main

import (
	"context"
	"fmt"
	"log"
	"os"

	paseto "github.com/murilo05/JobScheduler/internal/adapter/auth"
	"github.com/murilo05/JobScheduler/internal/adapter/config"
	"github.com/murilo05/JobScheduler/internal/adapter/handler/http"
	httpHandler "github.com/murilo05/JobScheduler/internal/adapter/handler/http"
	"github.com/murilo05/JobScheduler/internal/adapter/repository"
	"github.com/murilo05/JobScheduler/internal/adapter/repository/storage"
	"github.com/murilo05/JobScheduler/internal/core/service"
	"go.uber.org/zap"
)

func main() {

	config, err := config.New()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	zap, _ := zap.NewProduction()
	defer zap.Sync()
	logger := zap.Sugar()
	logger.Info("Starting the application",
		"app", config.App.Name,
		"env", config.App.Env,
	)

	ctx := context.Background()

	// Init token service
	token, err := paseto.New(config.Token)
	if err != nil {
		logger.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// Dependency injection
	// User
	storage, err := storage.NewStorage(ctx, config.DB, logger)
	userRepo := repository.NewRepository(storage, logger)
	userService := service.NewUserService(userRepo, logger)
	userHandler := httpHandler.NewUserHandler(userService)

	// Auth
	authService := service.NewAuthService(userRepo, token, logger)
	authHandler := httpHandler.NewAuthHandler(authService)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
		token,
		*userHandler,
		*authHandler,
	)
	if err != nil {
		logger.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	logger.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		logger.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
