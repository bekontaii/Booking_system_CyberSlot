package user

import "net/http"

// RegisterRoutes wires user handlers to the provided router.
func RegisterRoutes(mux *http.ServeMux) {
	repo := NewInMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/users", handler.HandleUsers)
	mux.HandleFunc("/users/", handler.HandleUserByID)
}
