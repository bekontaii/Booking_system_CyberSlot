package app

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/router"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/storage"
	"github.com/joho/godotenv"
)

// RunServer starts the HTTP server.
func RunServer() error {
	_ = godotenv.Load()
	addr := resolveAddr()

	db, err := storage.NewPostgres()
	if err != nil {
		log.Printf("PostgreSQL connection failed: %v", err)
		log.Println("Starting server in in-memory mode (fallback/demo)...")
		db = nil
	} else {
		defer db.Close()
		log.Println("Connected to PostgreSQL successfully")
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      router.New(db),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("server listening on %s", addr)
	return srv.ListenAndServe()
}

func resolveAddr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return ":8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}
