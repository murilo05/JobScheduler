package client

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/murilo05/JobScheduler/internal/adapter/config"
	"go.uber.org/zap"
)

type Postgres struct {
	db PG
	zap.SugaredLogger
}

type PG struct {
	*pgxpool.Pool
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

func NewDatabase(ctx context.Context, configPG *config.DB, logger *zap.SugaredLogger) *Postgres {
	pg, err := NewPostgres(ctx, configPG)
	if err != nil {
		logger.Error("Error initializing database connection: %s", err)
		os.Exit(1)
	}

	logger.Info("Successfully connected to the database", "db", configPG.Connection)

	return &Postgres{
		db:            *pg,
		SugaredLogger: *logger,
	}

}

func NewPostgres(ctx context.Context, config *config.DB) (*PG, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &PG{
		db,
		&psql,
		url,
	}, nil
}

func (db *PG) ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	} else {
		return fmt.Sprintf("unexpected error: %s", err)
	}
}

func (db *PG) Close() {
	db.Pool.Close()
}
