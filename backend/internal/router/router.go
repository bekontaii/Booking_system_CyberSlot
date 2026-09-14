package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	var clubRepo club.Repository
	var pcRepo pc.Repository
	var clubChecker booking.PCClubAvailabilityChecker
	if db != nil {
		userRepo = user.NewPostgresRepository(db)
		bookingRepo = booking.NewPostgresRepository(db)
		clubRepo = club.NewPostgresRepository(db)
		pcRepo = pc.NewPostgresRepository(db)
		clubChecker = booking.NewPostgresPCClubChecker(db)
	} else {
		userRepo = user.NewInMemoryRepository()
		bookingRepo = booking.NewInMemoryRepository()
		clubRepo = club.NewInMemoryRepository()
		pcRepo = pc.NewInMemoryRepository()
		clubChecker = booking.NewInMemoryPCClubChecker()
		seedInMemoryData(userRepo, clubRepo, pcRepo)
	}

	authService := auth.NewService(userRepo, secret, expireHours)
	authHandler := auth.NewHandler(authService)

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	bookingService := booking.NewService(bookingRepo, clubChecker, booking.DefaultExpiration)
	bookingHandler := booking.NewHandler(bookingService)

	clubService := club.NewService(clubRepo)
	clubHandler := club.NewHandler(clubService)

	pcService := pc.NewService(pcRepo)
	pcHandler := pc.NewHandler(pcService)

	publicMux := http.NewServeMux()
	staticDir := http.Dir(filepath.Join("web", "static"))
	publicMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Public auth routes (new + backward compatible)
	authMux := http.NewServeMux()
	auth.RegisterRoutes(authMux, authService)
	publicMux.Handle("/api/auth/", http.StripPrefix("/api/auth", authMux))
	publicMux.Handle("/auth/", http.StripPrefix("/auth", authMux))

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

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/logout", authHandler.Logout)
	apiMux.HandleFunc("/profile", userHandler.HandleProfile)

	// USER + SITE_ADMIN
	apiMux.Handle("/bookings", requireRolesForMethods(map[string][]string{
		http.MethodPost: []string{user.RoleUser, user.RoleSiteAdmin},
		http.MethodGet:  []string{user.RoleSiteAdmin},
	})(http.HandlerFunc(bookingHandler.HandleBookings)))
	apiMux.Handle("/bookings/club", middleware.RequireRoles(user.RoleClubAdmin, user.RoleSiteAdmin)(
		http.HandlerFunc(bookingHandler.HandleClubBookings),
	))
	apiMux.Handle("/bookings/", middleware.RequireRoles(user.RoleUser, user.RoleSiteAdmin)(
		requireBookingDeleteAccess(bookingService)(http.HandlerFunc(bookingHandler.HandleBookingByID)),
	))

	// USER (read) + SITE_ADMIN (full)
	apiMux.Handle("/clubs", requireRolesForMethods(map[string][]string{
		http.MethodPost:   []string{user.RoleSiteAdmin},
		http.MethodPut:    []string{user.RoleSiteAdmin},
		http.MethodPatch:  []string{user.RoleSiteAdmin},
		http.MethodDelete: []string{user.RoleSiteAdmin},
	})(http.HandlerFunc(clubHandler.HandleClubs)))
	apiMux.Handle("/clubs/", requireRolesForMethods(map[string][]string{
		http.MethodPut:    []string{user.RoleSiteAdmin},
		http.MethodPatch:  []string{user.RoleSiteAdmin},
		http.MethodDelete: []string{user.RoleSiteAdmin},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/activate") {
			clubHandler.HandleActivateClub(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/deactivate") {
			clubHandler.HandleDeactivateClub(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/pcs") {
			pcHandler.HandlePCsByClub(w, r)
			return
		}
		clubHandler.HandleClubByID(w, r)
	})))

	// CLUB_ADMIN (own club) + SITE_ADMIN
	pcAdminChain := requirePCClubScopeForMethods(pcRepo, []string{
		http.MethodPost, http.MethodPut, http.MethodDelete,
	})(requireRolesForMethods(map[string][]string{
		http.MethodPost:   []string{user.RoleClubAdmin, user.RoleSiteAdmin},
		http.MethodPut:    []string{user.RoleClubAdmin, user.RoleSiteAdmin},
		http.MethodDelete: []string{user.RoleClubAdmin, user.RoleSiteAdmin},
	})(http.HandlerFunc(pcHandler.HandlePCs)))
	apiMux.Handle("/pcs", pcAdminChain)

	pcByIDChain := requirePCClubScopeForMethods(pcRepo, []string{
		http.MethodPut, http.MethodDelete,
	})(requireRolesForMethods(map[string][]string{
		http.MethodPut:    []string{user.RoleClubAdmin, user.RoleSiteAdmin},
		http.MethodDelete: []string{user.RoleClubAdmin, user.RoleSiteAdmin},
	})(http.HandlerFunc(pcHandler.HandlePCByID)))
	apiMux.Handle("/pcs/", pcByIDChain)

	// SITE_ADMIN: user management + role assignment
	apiMux.Handle("/users", middleware.RequireRole(user.RoleSiteAdmin)(http.HandlerFunc(userHandler.HandleUsers)))
	apiMux.Handle("/users/", middleware.RequireRole(user.RoleSiteAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/role") {
			userHandler.HandleUserRoleByID(w, r)
			return
		}
		userHandler.HandleUserByID(w, r)
	})))

	protected := middleware.AuthMiddleware(secret)(apiMux)
	publicMux.Handle("/api/", http.StripPrefix("/api", protected))

	appDir := filepath.Join("web", "app")
	appIndex := filepath.Join(appDir, "index.html")
	if !fileExists(appIndex) {
		appDir = filepath.Join("backend", "web", "app")
		appIndex = filepath.Join(appDir, "index.html")
	}

	if fileExists(appIndex) {
		fileServer := http.FileServer(http.Dir(appDir))
		publicMux.Handle("/assets/", fileServer)
		publicMux.Handle("/logo.jpg", fileServer)
		publicMux.Handle("/favicon.ico", fileServer)
		publicMux.Handle("/", serveSPA(appIndex, appDir))
	} else {
		tmplPath := filepath.Join("web", "templates", "*.html")
		if matches, _ := filepath.Glob(tmplPath); len(matches) == 0 {
			tmplPath = filepath.Join("backend", "web", "templates", "*.html")
		}
		templates := template.Must(template.ParseGlob(tmplPath))
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

	return corsMiddleware(middleware.Logger(publicMux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func seedInMemoryData(userRepo user.Repository, clubRepo club.Repository, pcRepo pc.Repository) {
	adminHash, _ := auth.HashPassword("admin123")
	_, _ = userRepo.Create(user.User{
		Name:         "Admin",
		Surname:      "System",
		Email:        "admin@cyberslot.kz",
		Username:     "admin",
		Role:         user.RoleSiteAdmin,
		PasswordHash: adminHash,
	})

	playerHash, _ := auth.HashPassword("user123")
	_, _ = userRepo.Create(user.User{
		Name:         "Player",
		Surname:      "One",
		Email:        "player@cyberslot.kz",
		Username:     "player",
		Role:         user.RoleUser,
		PasswordHash: playerHash,
	})

	clubs := []club.Club{
		{Name: "Top Game", City: "Astana", Address: "Astana, Dinmukhamed Kunayev Street 23", IsActive: true},
		{Name: "BRO", City: "Astana", Address: "Astana, Kuishi Dina Street 31", IsActive: true},
		{Name: "Xan.exe", City: "Astana", Address: "Astana, Heydar Aliyev Street 3", IsActive: true},
		{Name: "Prime Game Hub", City: "Astana", Address: "Astana, Kerei and Zhanibek Khandar Street 14/2", IsActive: true},
		{Name: "Yamato Cyber Club", City: "Astana", Address: "Astana, Syganak Street 21/1", IsActive: true},
	}

	for _, c := range clubs {
		created, err := clubRepo.Create(c)
		if err != nil {
			continue
		}
		for i := 1; i <= 10; i++ {
			_, _ = pcRepo.Create(pc.PC{
				ClubID:   created.ID,
				PCNumber: i,
				Status:   pc.StatusActive,
			})
		}
	}
}

func serveSPA(indexPath, appDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(appDir))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			rel := strings.TrimPrefix(r.URL.Path, "/")
			rel = filepath.Clean(rel)
			if rel != "." && fileExists(filepath.Join(appDir, rel)) {
				fs.ServeHTTP(w, r)
				return
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

func requireRolesForMethods(methodRoles map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, restricted := methodRoles[r.Method]
			if !restricted {
				next.ServeHTTP(w, r)
				return
			}
			middleware.RequireRoles(roles...)(next).ServeHTTP(w, r)
		})
	}
}

func requirePCClubScopeForMethods(pcRepo pc.Repository, methods []string) func(http.Handler) http.Handler {
	methodSet := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		methodSet[m] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, restricted := methodSet[r.Method]; !restricted {
				next.ServeHTTP(w, r)
				return
			}

			authUser, ok := middleware.GetAuthUser(r.Context())
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if authUser.Role == user.RoleSiteAdmin {
				next.ServeHTTP(w, r)
				return
			}
			if authUser.Role != user.RoleClubAdmin || authUser.ClubID == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			targetClubID, ok := resolvePCClubID(r, pcRepo)
			if !ok || targetClubID != *authUser.ClubID {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func requireBookingDeleteAccess(bookingService *booking.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}

			authUser, ok := middleware.GetAuthUser(r.Context())
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if authUser.Role == user.RoleSiteAdmin {
				next.ServeHTTP(w, r)
				return
			}
			if authUser.Role != user.RoleUser {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			bookingID, ok := parseIDFromPrefix(r.URL.Path, "/bookings/")
			if !ok {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			entity, exists := bookingService.GetBookingByID(bookingID)
			if !exists {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			if entity.UserID != authUser.UserID {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func resolvePCClubID(r *http.Request, pcRepo pc.Repository) (int, bool) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return 0, false
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		var payload struct {
			ClubID int `json:"club_id"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return 0, false
		}
		return payload.ClubID, payload.ClubID > 0
	case http.MethodPut, http.MethodDelete:
		pcID, ok := parseIDFromPrefix(r.URL.Path, "/pcs/")
		if !ok {
			return 0, false
		}
		entity, err := pcRepo.GetByID(pcID)
		if err != nil {
			return 0, false
		}
		return entity.ClubID, entity.ClubID > 0
	default:
		return 0, false
	}
}

func parseIDFromPrefix(path, prefix string) (int, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	idPart := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if idPart == "" || strings.Contains(idPart, "/") {
		return 0, false
	}
	id, err := strconv.Atoi(idPart)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
