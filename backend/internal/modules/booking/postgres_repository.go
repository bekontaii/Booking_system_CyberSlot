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

func (r *PostgresRepository) Create(booking Booking) (Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.bookings (user_id, pc_id, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		booking.UserID,
		booking.PCID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
	).Scan(&booking.ID)

	return booking, err
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

func (r *PostgresRepository) GetByID(id int) (Booking, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var b Booking
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, pc_id, start_time, end_time, status
		FROM public.bookings
		WHERE id = $1
	`, id).Scan(&b.ID, &b.UserID, &b.PCID, &b.StartTime, &b.EndTime, &b.Status)
	if err != nil {
		return Booking{}, false
	}
	return b, true
}

func (r *PostgresRepository) GetByClubID(clubID int) []Booking {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.user_id, b.pc_id, b.start_time, b.end_time, b.status
		FROM public.bookings b
		JOIN public.pcs p ON p.id = b.pc_id
		WHERE p.club_id = $1
		ORDER BY b.id
	`, clubID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.PCID, &b.StartTime, &b.EndTime, &b.Status); err == nil {
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

func (r *PostgresRepository) Update(booking Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(
		ctx,
		`
			UPDATE public.bookings
			SET user_id = $1, pc_id = $2, start_time = $3, end_time = $4, status = $5
			WHERE id = $6
		`,
		booking.UserID,
		booking.PCID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.ID,
	)

	return err
}

func (r *PostgresRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(
		ctx,
		`DELETE FROM public.bookings WHERE id = $1`,
		id,
	)

	return err
}
