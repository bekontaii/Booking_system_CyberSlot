package router

import (
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

// New builds the application router with public and protected routes.
func New() http.Handler {
	templates := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))

	publicMux := http.NewServeMux()

	// Static assets
	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	secret, expireHours := resolveJWTConfig()
	authService := auth.NewService(secret, expireHours)

	// Auth API (public): /auth/login and /auth/register
	authMux := http.NewServeMux()
	auth.RegisterRoutes(authMux, authService)
	publicMux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	// Public HTML pages
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

	// Protected routes
	protectedMux := http.NewServeMux()

	// User module
	user.RegisterRoutes(protectedMux)

	// Club and PC modules (use sub-muxes to avoid /clubs/ route conflicts)
	clubMux := http.NewServeMux()
	club.RegisterRoutes(clubMux)

	pcMux := http.NewServeMux()
	pc.RegisterRoutes(pcMux)

	protectedMux.Handle("/clubs", clubMux)
	protectedMux.Handle("/clubs/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delegate /clubs/{id}/pcs to the PC module; others go to the club module.
		if strings.HasSuffix(r.URL.Path, "/pcs") || strings.HasSuffix(r.URL.Path, "/pcs/") {
			pcMux.ServeHTTP(w, r)
			return
		}
		clubMux.ServeHTTP(w, r)
	}))

	protectedMux.Handle("/pcs", pcMux)
	protectedMux.Handle("/pcs/", pcMux)

	// Booking module
	protectedMux.HandleFunc("/bookings", booking.HandleBookings)

	// Payment module (stub handler until implemented)
	protectedMux.HandleFunc("/payment", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("payment handler not implemented"))
	})

	// Everything else requires JWT
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
