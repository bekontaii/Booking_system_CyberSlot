package club

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	repo := NewInMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/clubs", handler.HandleClubs)
	mux.HandleFunc("/clubs/", handler.HandleClubByID)
}
