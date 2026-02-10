package router

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/middleware"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/auth"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/booking"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/club"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/pc"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
)

func New() http.Handler {
	templates := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))

	publicMux := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	secret, expireHours := resolveJWTConfig()
	authService := auth.NewService(secret, expireHours)
	authHandler := auth.NewHandler(authService)

	// Public HTML pages (no JWT)
	publicMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		renderTemplate(w, templates, "login.html")
	})

	publicMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		renderTemplate(w, templates, "register.html")
	})

	publicMux.HandleFunc("/clubs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		renderTemplate(w, templates, "clubs.html")
	})

	publicMux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		renderTemplate(w, templates, "booking.html")
	})

	publicMux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!DOCTYPE html><html><head><meta charset="utf-8"><title>Logout</title></head><body><script>
fetch('/api/logout', {method:'POST', headers:{'Authorization':'Bearer ' + (localStorage.getItem('jwt_token') || '')}})
  .finally(function(){ localStorage.removeItem('jwt_token'); window.location = '/login'; });
</script></body></html>`)
	})

	// Public auth API (no JWT): /auth/login, /auth/register
	publicAuthMux := http.NewServeMux()
	auth.RegisterRoutes(publicAuthMux, authService)
	publicMux.Handle("/auth/", http.StripPrefix("/auth", publicAuthMux))

	// Protected routes (JWT required)
	protectedMux := http.NewServeMux()

	// API routes under /api/*
	apiMux := http.NewServeMux()

	// Logout API (protected)
	apiMux.HandleFunc("/logout", authHandler.Logout)

	// User module
	user.RegisterRoutes(apiMux)

	// Club and PC modules (use sub-muxes to avoid /clubs/ route conflicts)
	clubMux := http.NewServeMux()
	club.RegisterRoutes(clubMux)

	pcMux := http.NewServeMux()
	pc.RegisterRoutes(pcMux)

	apiMux.Handle("/clubs", clubMux)
	apiMux.Handle("/clubs/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/pcs") || strings.HasSuffix(r.URL.Path, "/pcs/") {
			pcMux.ServeHTTP(w, r)
			return
		}
		clubMux.ServeHTTP(w, r)
	}))

	apiMux.Handle("/pcs", pcMux)
	apiMux.Handle("/pcs/", pcMux)

	// Booking module
	apiMux.HandleFunc("/bookings", booking.HandleBookings)

	// Payment module (stub handler until implemented)
	apiMux.HandleFunc("/payment", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("payment handler not implemented"))
	})

	protectedMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	publicMux.Handle("/", middleware.AuthMiddleware(secret)(protectedMux))

	return middleware.Logger(publicMux)
}

func renderTemplate(w http.ResponseWriter, templates *template.Template, name string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, nil); err != nil {
		http.Error(w, "template rendering error", http.StatusInternalServerError)
	}
}

func resolveJWTConfig() (string, int) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("JWT_SECRET not set; using development default")
		secret = "dev-secret"
	}

	expireHours := 24
	if value := strings.TrimSpace(os.Getenv("JWT_EXPIRE_HOURS")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			log.Println("JWT_EXPIRE_HOURS invalid; using default 24h")
		} else {
			expireHours = parsed
		}
	}

	return secret, expireHours
}
