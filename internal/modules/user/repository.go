package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			username,
			email,
			password_hash,
			name,
			surname
		) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Surname,
	)

	return err
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password_hash,
			name,
			surname,
			created_at,
			updated_at
		FROM users
		WHERE username = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Surname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
