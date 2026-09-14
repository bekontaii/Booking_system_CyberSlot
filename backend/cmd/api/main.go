package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/router"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/storage"
)

func main() {
	loadEnvFile(".env")
	loadEnvFile("backend/.env")

	addr := resolveAddr()
	ensureJWTSecret()

	db, err := storage.NewPostgres()
	if err != nil {
		log.Printf("PostgreSQL connection failed: %v", err)
		log.Println("Starting server in in-memory mode (fallback/demo)...")
		db = nil
	} else {
		defer db.Close()
		log.Println("Connected to PostgreSQL successfully")
	}

	handler := router.New(db)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("server listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func loadEnvFile(path string) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")

		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
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

func ensureJWTSecret() {
	if strings.TrimSpace(os.Getenv("JWT_SECRET")) != "" {
		return
	}

	log.Println("JWT_SECRET not set; using development default")
	_ = os.Setenv("JWT_SECRET", "dev-secret")
}
