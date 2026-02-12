package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/middleware"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/auth"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/booking"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/club"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/pc"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
)

func New(db *pgxpool.Pool) http.Handler {
	secret, expireHours := resolveJWTConfig()

	var userRepo user.Repository
	var bookingRepo booking.Repository
	if db != nil {
		userRepo = user.NewPostgresRepository(db)
		bookingRepo = booking.NewPostgresRepository(db)
	} else {
		userRepo = user.NewInMemoryRepository()
		bookingRepo = booking.NewInMemoryRepository()
	}

	// Shared repositories
	authService := auth.NewService(userRepo, secret, expireHours)
	authHandler := auth.NewHandler(authService)

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	bookingService := booking.NewService(bookingRepo, booking.DefaultExpiration)
	bookingHandler := booking.NewHandler(bookingService, userRepo)

	var clubRepo club.Repository
	var pcRepo pc.Repository
	if db != nil {
		clubRepo = club.NewPostgresRepository(db)
		pcRepo = pc.NewPostgresRepository(db)
	} else {
		clubRepo = club.NewInMemoryRepository()
		pcRepo = pc.NewInMemoryRepository()
	}
	clubService := club.NewService(clubRepo)
	clubHandler := club.NewHandler(clubService)

	pcService := pc.NewService(pcRepo)
	pcHandler := pc.NewHandler(pcService)

	publicMux := http.NewServeMux()

	// Static files for legacy templates
	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Public auth endpoints
	authMux := http.NewServeMux()
	auth.RegisterRoutes(authMux, authService)
	publicMux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	// Public logout page clears token on client
	publicMux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<script>
localStorage.removeItem('jwt_token');
location='/login';
</script>`)
	})

	// Protected API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/logout", authHandler.Logout)
	apiMux.HandleFunc("/users", userHandler.HandleUsers)
	apiMux.HandleFunc("/users/", userHandler.HandleUserByID)
	apiMux.HandleFunc("/profile", userHandler.HandleProfile)

	apiMux.HandleFunc("/bookings", bookingHandler.HandleBookings)
	apiMux.HandleFunc("/bookings/", bookingHandler.HandleBookingByID)

	apiMux.HandleFunc("/pcs", pcHandler.HandlePCs)
	apiMux.HandleFunc("/pcs/", pcHandler.HandlePCByID)

	apiMux.HandleFunc("/clubs", clubHandler.HandleClubs)
	apiMux.HandleFunc("/clubs/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/pcs") {
			pcHandler.HandlePCsByClub(w, r)
			return
		}
		clubHandler.HandleClubByID(w, r)
	})

	protected := middleware.AuthMiddleware(secret)(apiMux)
	publicMux.Handle("/api/", http.StripPrefix("/api", protected))

	appDir := filepath.Join("web", "app")
	appIndex := filepath.Join(appDir, "index.html")
	hasSPA := fileExists(appIndex)

	if hasSPA {
		fileServer := http.FileServer(http.Dir(appDir))
		publicMux.Handle("/assets/", fileServer)
		publicMux.Handle("/logo.jpg", fileServer)
		publicMux.Handle("/favicon.ico", fileServer)
		publicMux.Handle("/", serveSPA(appIndex, appDir))
	} else {
		templates := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))

		publicMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
		})
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
	}

	return middleware.Logger(publicMux)
}

func serveSPA(indexPath, appDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(appDir))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			rel := strings.TrimPrefix(r.URL.Path, "/")
			rel = filepath.Clean(rel)
			if rel != "." {
				if fileExists(filepath.Join(appDir, rel)) {
					fs.ServeHTTP(w, r)
					return
				}
			}
		}
		http.ServeFile(w, r, indexPath)
	}
}

func renderTemplate(w http.ResponseWriter, t *template.Template, name string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.ExecuteTemplate(w, name, nil)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
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
