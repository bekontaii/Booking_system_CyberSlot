package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/auth"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := user.NewRepository(db)

	authService := auth.NewService(userRepo)

	authHandler := auth.NewHandler(authService, jwtSecret)

	mux := http.NewServeMux()

	auth.RegisterRoutes(mux, authHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
