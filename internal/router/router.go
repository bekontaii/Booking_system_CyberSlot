package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/middleware"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/auth"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/booking"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/club"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/pc"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
)

func New(db *pgxpool.Pool) http.Handler {
	// ---------- templates ----------
	templates := template.Must(
		template.ParseGlob(filepath.Join("web", "templates", "*.html")),
	)

	publicMux := http.NewServeMux()

	// ---------- static ----------
	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(staticDir)),
	)

	// ---------- auth ----------
	secret, expireHours := resolveJWTConfig()
	authRepo := user.NewPostgresRepository(db)
	authService := auth.NewService(authRepo, secret, expireHours)
	authHandler := auth.NewHandler(authService)

	// ---------- public pages ----------
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

	// ---------- public auth api ----------
	authMux := http.NewServeMux()
	auth.RegisterRoutes(authMux, authService)
	publicMux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	// ---------- protected api ----------
	apiMux := http.NewServeMux()

	// logout
	apiMux.HandleFunc("/logout", authHandler.Logout)

	// ---------- USER (Postgres) ----------
	user.RegisterRoutes(apiMux, db)

	// ---------- BOOKING (Postgres) ----------
	bookingRepo := booking.NewPostgresRepository(db)
	bookingService := booking.NewService(bookingRepo, booking.DefaultExpiration)
	bookingHandler := booking.NewHandler(bookingService)

	apiMux.HandleFunc("/bookings", bookingHandler.HandleBookings)

	// ---------- CLUB / PC ----------
	clubMux := http.NewServeMux()
	club.RegisterRoutes(clubMux)

	pcMux := http.NewServeMux()
	pc.RegisterRoutes(pcMux)

	apiMux.Handle("/clubs/", clubMux)
	apiMux.Handle("/pcs/", pcMux)

	// ---------- middleware ----------
	protected := middleware.AuthMiddleware(secret)(apiMux)
	publicMux.Handle("/api/", http.StripPrefix("/api", protected))

	return middleware.Logger(publicMux)
}

/* ---------------- helpers ---------------- */

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
