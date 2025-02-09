package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/murilo05/JobScheduler/internal/adapter/storage/postgres"
	"github.com/murilo05/JobScheduler/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Insert("public.user").
		Columns("name", "email", "password", "role", "created_at", "updated_at").
		Values(user.Name, user.Email, user.Password, "customer", time.Now(), time.Now()).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("public.user").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("public.user").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	name := nullString(user.Name)
	email := nullString(user.Email)
	password := nullString(user.Password)
	role := nullString(string(user.Role))

	query := ur.db.QueryBuilder.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("email", sq.Expr("COALESCE(?, email)", email)).
		Set("password", sq.Expr("COALESCE(?, password)", password)).
		Set("role", sq.Expr("COALESCE(?, role)", role)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := ur.db.QueryBuilder.Delete("public.user").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Print("aa", err)
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		fmt.Print("aaxxx", err)
		return err
	}

	return nil
}
