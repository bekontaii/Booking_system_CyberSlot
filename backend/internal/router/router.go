package router

import (
	"fmt"
	middleware2 "github.com/bekontaii/Booking_system_CyberSlot/backend/internal/middleware"
	auth2 "github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/auth"
	booking2 "github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/booking"
	"github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/club"
	"github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/pc"
	user2 "github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/user"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) http.Handler {
	templates := template.Must(
		template.ParseGlob(filepath.Join("web", "templates", "*.html")),
	)

	publicMux := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(staticDir)),
	)

	secret, expireHours := resolveJWTConfig()
	authRepo := user2.NewPostgresRepository(db)
	authService := auth2.NewService(authRepo, secret, expireHours)
	authHandler := auth2.NewHandler(authService)

	publicMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "login.html")
	})

	publicMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "register.html")
	})

	publicMux.HandleFunc("/clubs", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "clubs.html")
	})

	publicMux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "booking.html")
	})

	publicMux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<script>
fetch('/api/logout',{method:'POST',headers:{'Authorization':'Bearer '+localStorage.getItem('jwt_token')}})
.finally(()=>{localStorage.removeItem('jwt_token');location='/login';});
</script>`)
	})

	authMux := http.NewServeMux()
	auth2.RegisterRoutes(authMux, authService)
	publicMux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	publicMux.HandleFunc("/api/logout", authHandler.Logout)

	apiMux := http.NewServeMux()

	// logout
	apiMux.HandleFunc("/logout", authHandler.Logout)

	user2.RegisterRoutes(apiMux, db)

	bookingRepo := booking2.NewPostgresRepository(db)
	bookingService := booking2.NewService(bookingRepo, booking2.DefaultExpiration)
	bookingHandler := booking2.NewHandler(bookingService)

	apiMux.HandleFunc("/bookings", bookingHandler.HandleBookings)
	apiMux.HandleFunc("/bookings/", bookingHandler.HandleBookingByID)

	clubMux := http.NewServeMux()
	club.RegisterRoutes(clubMux)

	pcMux := http.NewServeMux()
	pc.RegisterRoutes(pcMux)

	apiMux.Handle("/clubs/", clubMux)
	apiMux.Handle("/pcs/", pcMux)

	protected := middleware2.AuthMiddleware(secret)(apiMux)
	publicMux.Handle("/api/", http.StripPrefix("/api", protected))

	return middleware2.Logger(publicMux)
}

func renderTemplate(w http.ResponseWriter, t *template.Template, name string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.ExecuteTemplate(w, name, nil)
}

func resolveJWTConfig() (string, int) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	expire := 24
	if v := os.Getenv("JWT_EXPIRE_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			expire = n
		}
	}

	return secret, expire
}
