package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *pgxpool.Pool) {
	repo := NewPostgresRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/users", handler.HandleUsers)
	mux.HandleFunc("/users/", handler.HandleUserByID)
}
