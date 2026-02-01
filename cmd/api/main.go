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

	// 1️⃣ JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// 2️⃣ Database connection
	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3️⃣ Repository
	userRepo := user.NewRepository(db)

	// 4️⃣ Service
	authService := auth.NewService(userRepo)

	// 5️⃣ Handler
	authHandler := auth.NewHandler(authService, jwtSecret)

	// 6️⃣ Router (ServeMux)
	mux := http.NewServeMux()

	// 7️⃣ Register routes
	auth.RegisterRoutes(mux, authHandler)

	// 8️⃣ Start server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
