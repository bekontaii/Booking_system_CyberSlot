package pc

import "net/http"

// RegisterRoutes wires PC handlers to the provided router.
func RegisterRoutes(mux *http.ServeMux) {
	repo := NewInMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/pcs", handler.HandlePCs)
	mux.HandleFunc("/pcs/", handler.HandlePCByID)
	mux.HandleFunc("/clubs/", handler.HandlePCsByClub)
}
