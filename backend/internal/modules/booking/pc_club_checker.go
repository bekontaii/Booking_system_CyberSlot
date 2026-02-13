package booking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPCClubChecker struct {
	db *pgxpool.Pool
}

func NewPostgresPCClubChecker(db *pgxpool.Pool) *PostgresPCClubChecker {
	return &PostgresPCClubChecker{db: db}
}

func (c *PostgresPCClubChecker) IsPCBookable(pcID int) bool {
	if pcID <= 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var isActive bool
	err := c.db.QueryRow(ctx, `
		SELECT cl.is_active
		FROM public.pcs p
		JOIN public.clubs cl ON cl.id = p.club_id
		WHERE p.id = $1
	`, pcID).Scan(&isActive)
	if err != nil {
		return false
	}
	return isActive
}

type InMemoryPCClubChecker struct{}

func NewInMemoryPCClubChecker() *InMemoryPCClubChecker {
	return &InMemoryPCClubChecker{}
}

func (c *InMemoryPCClubChecker) IsPCBookable(pcID int) bool {
	return pcID > 0
}

