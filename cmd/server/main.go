package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"usl-server/internal/auth"
	"usl-server/internal/config"
	"usl-server/internal/handlers"
	"usl-server/internal/logger"
	"usl-server/internal/middleware"
	"usl-server/internal/repositories"
	"usl-server/internal/services"
	usl "usl-server/internal/usl"
	uslHandlers "usl-server/internal/usl/handlers"

	"github.com/supabase-community/supabase-go"
)

type ApplicationContext struct {
	Config           *config.Config
	UserRepo         *repositories.UserRepository
	TrackerRepo      *repositories.TrackerRepository
	GuildRepo        *repositories.GuildRepository
	TrueSkillService *services.UserTrueSkillService
	Templates        *template.Template
	Logger           *slog.Logger
	Auth             *auth.DiscordAuth
}

func main() {
	// Setup logging first - write to both stdout AND log file
	appLogger := logger.SetupLogger(slog.LevelDebug, "logs/server.log")
	appLogger.Info("Starting USL server application")

	dependencies := initializeApplication(appLogger)
	server := setupHTTPServer(dependencies)
	startServer(server, dependencies.Config, dependencies.Logger)
}

func initializeApplication(appLogger *slog.Logger) *ApplicationContext {
	configuration := loadConfiguration(appLogger)
	supabaseClient := createSupabaseClient(configuration, appLogger)
	repositories := setupRepositories(supabaseClient, configuration, appLogger)
	services := setupServices(configuration, repositories, appLogger)
	templates := loadTemplates(appLogger)

	// Create unified Discord OAuth auth system
	discordAuth := auth.NewDiscordAuth(supabaseClient, configuration.USL.AdminDiscordIDs,
		configuration.Supabase.URL, configuration.Supabase.PublicURL, configuration.Supabase.AnonKey)

	return &ApplicationContext{
		Config:           configuration,
		UserRepo:         repositories.UserRepo,
		TrackerRepo:      repositories.TrackerRepo,
		GuildRepo:        repositories.GuildRepo,
		TrueSkillService: services,
		Templates:        templates,
		Logger:           appLogger,
		Auth:             discordAuth,
	}
}

func loadConfiguration(logger *slog.Logger) *config.Config {
	logger.Info("Loading configuration")
	configuration, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully")
	return configuration
}

func createSupabaseClient(configuration *config.Config, logger *slog.Logger) *supabase.Client {
	logger.Info("Initializing Supabase client", "url", configuration.Supabase.URL)
	client, err := supabase.NewClient(
		configuration.Supabase.URL,
		configuration.Supabase.ServiceRoleKey,
		nil,
	)
	if err != nil {
		logger.Error("Failed to initialize Supabase client", "error", err)
		os.Exit(1)
	}
	logger.Info("Supabase client initialized successfully")
	return client
}

type RepositoryCollection struct {
	UserRepo    *repositories.UserRepository
	TrackerRepo *repositories.TrackerRepository
	GuildRepo   *repositories.GuildRepository
}

func setupRepositories(client *supabase.Client, config *config.Config, logger *slog.Logger) *RepositoryCollection {
	logger.Info("Setting up repositories")
	return &RepositoryCollection{
		UserRepo:    repositories.NewUserRepository(client, config),
		TrackerRepo: repositories.NewTrackerRepository(client, config),
		GuildRepo:   repositories.NewGuildRepository(client, config),
	}
}

func setupServices(config *config.Config, repos *RepositoryCollection, logger *slog.Logger) *services.UserTrueSkillService {
	logger.Info("Setting up services")
	percentileConverter := services.NewPercentileConverter(config)
	mmrCalculator := services.NewMMRCalculator(config, percentileConverter)
	uncertaintyCalculator := services.NewEnhancedUncertaintyCalculator(config)
	dataTransformationService := services.NewDataTransformationService()

	return services.NewUserTrueSkillService(
		repos.TrackerRepo,
		repos.UserRepo,
		mmrCalculator,
		uncertaintyCalculator,
		dataTransformationService,
		config,
	)
}

func createTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...any) map[string]any {
			dict := make(map[string]any)
			for i := 0; i < len(values); i += 2 {
				if i+1 < len(values) {
					dict[values[i].(string)] = values[i+1]
				}
			}
			return dict
		},
		"slice": func(values ...any) []any {
			return values
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"printf": func(format string, args ...any) string {
			return fmt.Sprintf(format, args...)
		},
		"lt": func(a, b float64) bool {
			return a < b
		},
		"substr": func(s string, start, length int) string {
			if start >= len(s) {
				return ""
			}
			end := min(start+length, len(s))
			return s[start:end]
		},
	}
}

func loadTemplates(logger *slog.Logger) *template.Template {
	logger.Info("Loading templates")

	tmpl := template.New("app")
	tmpl = tmpl.Funcs(createTemplateFunctions())

	// Parse all template files
	patterns := []string{
		"templates/*.html",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			logger.Error("Failed to glob templates", "pattern", pattern, "error", err)
			os.Exit(1)
		}

		if len(matches) > 0 {
			logger.Info("Loading templates", "pattern", pattern, "count", len(matches))
			tmpl = template.Must(tmpl.ParseFiles(matches...))
		}
	}

	logger.Info("Templates loaded successfully")
	return tmpl
}

func setupHTTPServer(deps *ApplicationContext) http.Handler {
	mux := http.NewServeMux()

	setupStaticRoutes(mux)
	setupHealthRoute(mux)      // Health check endpoint
	setupAuthRoutes(mux, deps) // Unified Discord OAuth routes
	setupHomeRoute(mux)
	setupGuildAwareRoutes(mux, deps) // New guild-aware routes
	setupUserRoutes(mux, deps)       // Legacy routes
	setupTrackerRoutes(mux, deps)    // Legacy routes
	setupTrueSkillRoutes(mux, deps)  // Legacy routes
	setupAPIRoutes(mux, deps)
	setupUSLRoutes(mux, deps) // USL-specific temporary migration routes

	// Create guild context middleware
	guildMiddleware := middleware.NewGuildContextMiddleware(deps.GuildRepo, deps.Logger)

	// Apply middleware chain
	handler := middleware.LoggingMiddleware(deps.Logger)(mux)
	// Apply guild context middleware globally now that it handles missing schema gracefully
	handler = guildMiddleware.GuildContext()(handler)

	return handler
}

func setupStaticRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}

func setupHealthRoute(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"healthy","service":"usl-server"}`)); err != nil {
			log.Printf("Failed to write health response: %v", err)
		}
	})
}

func setupGuildAwareRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	// Create new guild-aware handlers for future guild support
	// Currently no routes registered here to avoid conflicts with USL routes
	// TODO: Add support for dynamic guild slugs when needed
}

func setupAuthRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	// Redirect main app login to USL login for now
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
	})
	mux.HandleFunc("/usl/login", deps.Auth.LoginForm)
	mux.HandleFunc("/auth/callback", deps.Auth.AuthCallback)
	mux.HandleFunc("/auth/process", deps.Auth.ProcessTokens)
	mux.HandleFunc("/logout", deps.Auth.Logout)
	mux.HandleFunc("/usl/logout", deps.Auth.Logout) // USL uses same logout
}

func setupHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", handleHomeAndNotFound)
}

func handleHomeAndNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		render404Page(w)
		return
	}
	http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
}

func render404Page(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head><title>404 - Page Not Found</title></head>
<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
<h1>404 - Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<a href="/" style="color: #007cba;">← Go Home</a> | 
<a href="/usl/admin" style="color: #007cba;">USL Admin</a>
</body>
</html>`

	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write 404 response: %v", err)
	}
}

func setupUserRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	// Legacy routes - redirect to guild-aware routes
	mux.HandleFunc("/users", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/new", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/new", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/create", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/create", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/edit", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/edit", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/update", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/update", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/delete", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/delete", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/search", deps.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/search", http.StatusMovedPermanently)
	}))
}

func setupTrackerRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.TrueSkillService, deps.Templates)

	// Main app tracker routes - protected by unified Discord OAuth
	mux.HandleFunc("/trackers", deps.Auth.RequireAuth(trackerHandler.ListTrackers))
	mux.HandleFunc("/trackers/new", deps.Auth.RequireAuth(trackerHandler.NewTrackerForm))
	mux.HandleFunc("/trackers/create", deps.Auth.RequireAuth(trackerHandler.CreateTracker))
	mux.HandleFunc("/trackers/edit", deps.Auth.RequireAuth(trackerHandler.EditTrackerForm))
	mux.HandleFunc("/trackers/update", deps.Auth.RequireAuth(trackerHandler.UpdateTracker))
	mux.HandleFunc("/trackers/delete", deps.Auth.RequireAuth(trackerHandler.DeleteTracker))
	mux.HandleFunc("/trackers/search", deps.Auth.RequireAuth(trackerHandler.SearchTrackers))
}

func setupTrueSkillRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	trueskillHandler := handlers.NewTrueSkillHandler(deps.TrueSkillService, deps.Templates)

	// Main app TrueSkill routes - protected by unified Discord OAuth
	mux.HandleFunc("/trueskill/update-all", deps.Auth.RequireAuth(trueskillHandler.UpdateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/update-user", deps.Auth.RequireAuth(trueskillHandler.UpdateUserTrueSkill))
	mux.HandleFunc("/trueskill/recalculate", deps.Auth.RequireAuth(trueskillHandler.RecalculateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/stats", deps.Auth.RequireAuth(trueskillHandler.GetTrueSkillStats))
}

func setupAPIRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	// v1 API handlers (legacy)
	userHandler := handlers.NewUserHandler(deps.UserRepo, deps.Templates)
	trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.TrueSkillService, deps.Templates)
	trueskillHandler := handlers.NewTrueSkillHandler(deps.TrueSkillService, deps.Templates)

	// v2 API handlers (modern with pagination, filtering, bulk operations)
	v2UsersHandler := uslHandlers.NewV2UsersHandler(deps.UserRepo)
	v2TrackersHandler := uslHandlers.NewV2TrackersHandler(deps.TrackerRepo)

	// v1 API routes - protected by unified Discord OAuth (legacy)
	mux.HandleFunc("/api/users", deps.Auth.RequireAuth(userHandler.ListUsersAPI))
	mux.HandleFunc("/api/trackers", deps.Auth.RequireAuth(trackerHandler.ListTrackersAPI))
	mux.HandleFunc("/api/trueskill/update-all", deps.Auth.RequireAuth(trueskillHandler.UpdateAllUserTrueSkillAPI))

	// v2 API routes - modern paginated APIs (protected by unified Discord OAuth)
	mux.HandleFunc("/api/v2/users", deps.Auth.RequireAuth(v2UsersHandler.HandleUsers))
	mux.HandleFunc("/api/v2/users/bulk", deps.Auth.RequireAuth(v2UsersHandler.HandleUsersBulk))
	mux.HandleFunc("/api/v2/trackers", deps.Auth.RequireAuth(v2TrackersHandler.HandleTrackers))
	mux.HandleFunc("/api/v2/trackers/bulk", deps.Auth.RequireAuth(v2TrackersHandler.HandleTrackersBulk))
}

func setupUSLRoutes(mux *http.ServeMux, deps *ApplicationContext) {
	// TEMPORARY: USL migration handler - will be deleted after migration
	// Create dedicated Supabase client for USL operations
	supabaseClient, err := supabase.NewClient(
		deps.Config.Supabase.URL,
		deps.Config.Supabase.ServiceRoleKey,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create USL Supabase client: %v", err)
	}

	uslRepo := usl.NewUSLRepository(supabaseClient, deps.Config, deps.Logger)
	// NOTE: USL handlers no longer need their own auth - they use the unified auth
	uslHandler := uslHandlers.NewMigrationHandler(uslRepo, deps.Templates, deps.TrueSkillService, deps.Config)

	// USL Main Routes (redirect to login or admin based on auth status)
	mux.HandleFunc("/usl/", func(w http.ResponseWriter, r *http.Request) {
		if deps.Auth.IsAuthenticated(r) {
			http.Redirect(w, r, "/usl/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
		}
	})

	// USL User Management Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/users", deps.Auth.RequireAuth(uslHandler.ListUsers))
	mux.HandleFunc("/usl/users/search", deps.Auth.RequireAuth(uslHandler.SearchUsers))
	mux.HandleFunc("/usl/users/detail", deps.Auth.RequireAuth(uslHandler.UserDetail))
	mux.HandleFunc("/usl/users/new", deps.Auth.RequireAuth(uslHandler.NewUserForm))
	mux.HandleFunc("/usl/users/create", deps.Auth.RequireAuth(uslHandler.CreateUser))
	mux.HandleFunc("/usl/users/edit", deps.Auth.RequireAuth(uslHandler.EditUserForm))
	mux.HandleFunc("/usl/users/update", deps.Auth.RequireAuth(uslHandler.UpdateUser))
	mux.HandleFunc("/usl/users/delete", deps.Auth.RequireAuth(uslHandler.DeleteUser))
	mux.HandleFunc("/usl/users/update-trueskill", deps.Auth.RequireAuth(uslHandler.UpdateUserTrueSkill))

	// USL Tracker Management Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/trackers", deps.Auth.RequireAuth(uslHandler.ListTrackers))
	mux.HandleFunc("/usl/trackers/search", deps.Auth.RequireAuth(uslHandler.SearchTrackers))
	mux.HandleFunc("/usl/trackers/detail", deps.Auth.RequireAuth(uslHandler.TrackerDetail))
	mux.HandleFunc("/usl/trackers/new", deps.Auth.RequireAuth(uslHandler.NewTrackerForm))
	mux.HandleFunc("/usl/trackers/create", deps.Auth.RequireAuth(uslHandler.CreateTracker))
	mux.HandleFunc("/usl/trackers/edit", deps.Auth.RequireAuth(uslHandler.EditTrackerForm))
	mux.HandleFunc("/usl/trackers/update", deps.Auth.RequireAuth(uslHandler.UpdateTracker))

	// USL Admin Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/admin", deps.Auth.RequireAuth(uslHandler.AdminDashboard))

	// USL API Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/api/users", deps.Auth.RequireAuth(uslHandler.ListUsersAPI))
	mux.HandleFunc("/usl/api/trackers", deps.Auth.RequireAuth(uslHandler.ListTrackersAPI))
	mux.HandleFunc("/usl/api/leaderboard", deps.Auth.RequireAuth(uslHandler.GetLeaderboardAPI))
}

func startServer(server http.Handler, config *config.Config, logger *slog.Logger) {
	serverAddress := config.Server.Host + ":" + config.Server.Port

	logger.Info("Starting USL server",
		"address", serverAddress,
		"supabase_url", config.Supabase.URL)

	if err := http.ListenAndServe(serverAddress, server); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
