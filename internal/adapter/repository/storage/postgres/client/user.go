package client

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/murilo05/JobScheduler/internal/adapter/repository/storage"

	"github.com/murilo05/JobScheduler/internal/core/domain"
)

type Postgres struct {
	db *storage.PG
	_  storage.UserStorage
}

var _ storage.UserStorage = &Postgres{}

func (pg *Postgres) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := pg.db.QueryBuilder.Insert("public.user").
		Columns("name", "email", "password", "role", "created_at", "updated_at").
		Values(user.Name, user.Email, user.Password, "customer", time.Now(), time.Now()).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pg.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := pg.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (pg *Postgres) Get(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User

	query := pg.db.QueryBuilder.Select("*").
		From("public.user").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pg.db.QueryRow(ctx, sql, args...).Scan(
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

func (pg *Postgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := pg.db.QueryBuilder.Select("*").
		From("public.user").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pg.db.QueryRow(ctx, sql, args...).Scan(
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

func (pg *Postgres) List(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var user domain.User
	var users []domain.User

	query := pg.db.QueryBuilder.Select("*").
		From("public.user").
		OrderBy("id").
		Limit(limit).
		Offset(skip * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pg.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (pg *Postgres) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	name := nullString(user.Name)
	email := nullString(user.Email)
	password := nullString(user.Password)
	role := nullString(string(user.Role))

	query := pg.db.QueryBuilder.Update("public.user").
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

	err = pg.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := pg.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (pg *Postgres) Delete(ctx context.Context, id uint64) error {
	query := pg.db.QueryBuilder.Delete("public.user").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
