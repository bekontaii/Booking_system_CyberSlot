package pc

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	repo := NewInMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/pcs", handler.HandlePCs)
	mux.HandleFunc("/pcs/", handler.HandlePCByID)
	mux.HandleFunc("/clubs/", handler.HandlePCsByClub)
}
