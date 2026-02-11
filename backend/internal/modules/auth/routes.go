package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service *Service) {
	handler := NewHandler(service)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)
}
