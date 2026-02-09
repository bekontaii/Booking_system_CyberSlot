package auth

import "net/http"

// RegisterRoutes wires auth handlers to the provided router.
func RegisterRoutes(mux *http.ServeMux, service *Service) {
	handler := NewHandler(service)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)
}
