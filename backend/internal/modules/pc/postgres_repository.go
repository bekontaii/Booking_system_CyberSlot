package pc

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

func (r *PostgresRepository) Create(pc PC) (PC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO public.pcs (club_id, pc_number, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		pc.ClubID,
		pc.PCNumber,
		pc.Status,
	).Scan(&pc.ID, &pc.Created)
	return pc, err
}

func (r *PostgresRepository) GetAll() ([]PC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, club_id, pc_number, status, created_at
		FROM public.pcs
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcs := make([]PC, 0)
	for rows.Next() {
		var p PC
		if err := rows.Scan(&p.ID, &p.ClubID, &p.PCNumber, &p.Status, &p.Created); err != nil {
			return nil, err
		}
		pcs = append(pcs, p)
	}

	return pcs, nil
}

func (r *PostgresRepository) GetByID(id int) (PC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p PC
	err := r.db.QueryRow(
		ctx,
		`SELECT id, club_id, pc_number, status, created_at FROM public.pcs WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.ClubID, &p.PCNumber, &p.Status, &p.Created)
	if err != nil {
		return PC{}, ErrPCNotFound
	}

	return p, nil
}

func (r *PostgresRepository) GetByClubID(clubID int) ([]PC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, club_id, pc_number, status, created_at
		FROM public.pcs
		WHERE club_id = $1
		ORDER BY id
	`, clubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcs := make([]PC, 0)
	for rows.Next() {
		var p PC
		if err := rows.Scan(&p.ID, &p.ClubID, &p.PCNumber, &p.Status, &p.Created); err != nil {
			return nil, err
		}
		pcs = append(pcs, p)
	}

	return pcs, nil
}

func (r *PostgresRepository) Update(id int, pc PC) (PC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(
		ctx,
		`UPDATE public.pcs SET club_id = $1, pc_number = $2, status = $3 WHERE id = $4`,
		pc.ClubID,
		pc.PCNumber,
		pc.Status,
		id,
	)
	if err != nil {
		return PC{}, err
	}
	if cmd.RowsAffected() == 0 {
		return PC{}, ErrPCNotFound
	}

	pc.ID = id
	return pc, nil
}

func (r *PostgresRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := r.db.Exec(ctx, `DELETE FROM public.pcs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrPCNotFound
	}

	return nil
}
