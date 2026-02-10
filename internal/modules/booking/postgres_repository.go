package booking

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

func (r *PostgresRepository) Create(booking Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.bookings (user_id, pc_id, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		booking.UserID,
		booking.PCID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
	)

	return err
}

func (r *PostgresRepository) GetAll() []Booking {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, pc_id, start_time, end_time, status
		FROM public.bookings
		ORDER BY id
	`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.PCID,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
		); err == nil {
			bookings = append(bookings, b)
		}
	}

	return bookings
}

func (r *PostgresRepository) UpdateStatus(id int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(
		ctx,
		`UPDATE public.bookings SET status = $1 WHERE id = $2`,
		status,
		id,
	)

	return err
}
