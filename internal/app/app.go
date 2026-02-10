package app

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/router"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/storage"
)

func RunServer() error {
	addr := resolveAddr()

	db, err := storage.NewPostgres()
	if err != nil {
		return err
	}
	defer db.Close()

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
