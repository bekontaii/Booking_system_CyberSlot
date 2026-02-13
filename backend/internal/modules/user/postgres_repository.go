package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(user User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.users (name, surname, email, username, password_hash, role, club_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.ClubID,
	).Scan(&user.ID, &user.CreatedAt)

	return user, err
}

func (r *PostgresRepository) GetAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, name, surname, email, username, role, club_id, created_at
		FROM public.users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Surname,
			&u.Email,
			&u.Username,
			&u.Role,
			&u.ClubID,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
func (r *PostgresRepository) GetByID(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u User
	query := `
		SELECT id, name, surname, email, username, role, club_id, created_at
		FROM public.users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Surname,
		&u.Email,
		&u.Username,
		&u.Role,
		&u.ClubID,
		&u.CreatedAt,
	)

	if err != nil {
		return User{}, ErrUserNotFound
	}

	return u, nil
}

func (r *PostgresRepository) GetByUsername(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u User
	query := `
		SELECT id, name, surname, email, username, role, club_id, password_hash, created_at
		FROM public.users
		WHERE username = $1
	`

	err := r.db.QueryRow(ctx, query, username).Scan(
		&u.ID,
		&u.Name,
		&u.Surname,
		&u.Email,
		&u.Username,
		&u.Role,
		&u.ClubID,
		&u.PasswordHash,
		&u.CreatedAt,
	)
	if err != nil {
		return User{}, ErrUserNotFound
	}

	return u, nil
}

func (r *PostgresRepository) Update(id int, user User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE public.users
		SET name = $1, surname = $2, email = $3, username = $4, role = $5, club_id = $6
		WHERE id = $7
	`

	cmd, err := r.db.Exec(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Email,
		user.Username,
		user.Role,
		user.ClubID,
		id,
	)
	if err != nil {
		return User{}, err
	}

	if cmd.RowsAffected() == 0 {
		return User{}, ErrUserNotFound
	}

	user.ID = id
	return user, nil
}

func (r *PostgresRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(
		ctx,
		"DELETE FROM public.users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
