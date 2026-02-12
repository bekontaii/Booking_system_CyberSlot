package club

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

func (r *PostgresRepository) Create(club Club) (Club, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.clubs (name, city, address)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query, club.Name, club.City, club.Address).Scan(&club.ID)
	return club, err
}

func (r *PostgresRepository) GetAll() ([]Club, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, name, city, address
		FROM public.clubs
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clubs := make([]Club, 0)
	for rows.Next() {
		var c Club
		if err := rows.Scan(&c.ID, &c.Name, &c.City, &c.Address); err != nil {
			return nil, err
		}
		clubs = append(clubs, c)
	}

	return clubs, nil
}

func (r *PostgresRepository) GetByID(id int) (Club, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c Club
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, city, address FROM public.clubs WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.City, &c.Address)
	if err != nil {
		return Club{}, ErrClubNotFound
	}

	return c, nil
}

func (r *PostgresRepository) Update(id int, club Club) (Club, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(
		ctx,
		`UPDATE public.clubs SET name = $1, city = $2, address = $3 WHERE id = $4`,
		club.Name,
		club.City,
		club.Address,
		id,
	)
	if err != nil {
		return Club{}, err
	}
	if cmd.RowsAffected() == 0 {
		return Club{}, ErrClubNotFound
	}

	club.ID = id
	return club, nil
}

func (r *PostgresRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(ctx, `DELETE FROM public.clubs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrClubNotFound
	}

	return nil
}

