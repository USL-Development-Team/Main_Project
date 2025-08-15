package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
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

type AppDependencies struct {
	Config           *config.Config
	UserRepo         *repositories.UserRepository
	TrackerRepo      *repositories.TrackerRepository
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

func initializeApplication(appLogger *slog.Logger) *AppDependencies {
	configuration := loadConfiguration(appLogger)
	supabaseClient := createSupabaseClient(configuration, appLogger)
	repositories := setupRepositories(supabaseClient, configuration, appLogger)
	services := setupServices(configuration, repositories, appLogger)
	templates := loadTemplates(appLogger)

	// Create unified Discord OAuth auth system
	discordAuth := auth.NewDiscordAuth(supabaseClient, configuration.USL.AdminDiscordIDs,
		configuration.Supabase.URL, configuration.Supabase.PublicURL, configuration.Supabase.AnonKey)

	return &AppDependencies{
		Config:           configuration,
		UserRepo:         repositories.UserRepo,
		TrackerRepo:      repositories.TrackerRepo,
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
}

func setupRepositories(client *supabase.Client, config *config.Config, logger *slog.Logger) *RepositoryCollection {
	logger.Info("Setting up repositories")
	return &RepositoryCollection{
		UserRepo:    repositories.NewUserRepository(client, config),
		TrackerRepo: repositories.NewTrackerRepository(client, config),
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

func loadTemplates(logger *slog.Logger) *template.Template {
	logger.Info("Loading templates")
	return template.Must(template.ParseGlob("templates/*.html"))
}

func setupHTTPServer(deps *AppDependencies) http.Handler {
	mux := http.NewServeMux()

	setupStaticRoutes(mux)
	setupAuthRoutes(mux, deps) // Unified Discord OAuth routes
	setupHomeRoute(mux)
	setupUserRoutes(mux, deps)
	setupTrackerRoutes(mux, deps)
	setupTrueSkillRoutes(mux, deps)
	setupAPIRoutes(mux, deps)
	setupUSLRoutes(mux, deps) // USL-specific temporary migration routes

	// Wrap with logging middleware
	loggedHandler := middleware.LoggingMiddleware(deps.Logger)(mux)

	return loggedHandler
}

func setupStaticRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}

func setupAuthRoutes(mux *http.ServeMux, deps *AppDependencies) {
	// Unified Discord OAuth routes for both main app and USL
	mux.HandleFunc("/login", deps.Auth.LoginForm)
	mux.HandleFunc("/usl/login", deps.Auth.LoginForm) // USL uses same login
	mux.HandleFunc("/auth/callback", deps.Auth.AuthCallback)
	mux.HandleFunc("/auth/process", deps.Auth.ProcessTokens)
	mux.HandleFunc("/logout", deps.Auth.Logout)
	mux.HandleFunc("/usl/logout", deps.Auth.Logout) // USL uses same logout
}

func setupHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// This is an unmatched route, show custom 404
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			// Simple inline 404 for now - could use template later
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
			w.Write([]byte(html))
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	})
}

func setupUserRoutes(mux *http.ServeMux, deps *AppDependencies) {
	userHandler := handlers.NewUserHandler(deps.UserRepo, deps.Templates)

	// Main app user routes - protected by unified Discord OAuth
	mux.HandleFunc("/users", deps.Auth.RequireAuth(userHandler.ListUsers))
	mux.HandleFunc("/users/new", deps.Auth.RequireAuth(userHandler.NewUserForm))
	mux.HandleFunc("/users/create", deps.Auth.RequireAuth(userHandler.CreateUser))
	mux.HandleFunc("/users/edit", deps.Auth.RequireAuth(userHandler.EditUserForm))
	mux.HandleFunc("/users/update", deps.Auth.RequireAuth(userHandler.UpdateUser))
	mux.HandleFunc("/users/delete", deps.Auth.RequireAuth(userHandler.DeleteUser))
	mux.HandleFunc("/users/search", deps.Auth.RequireAuth(userHandler.SearchUsers))
}

func setupTrackerRoutes(mux *http.ServeMux, deps *AppDependencies) {
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

func setupTrueSkillRoutes(mux *http.ServeMux, deps *AppDependencies) {
	trueskillHandler := handlers.NewTrueSkillHandler(deps.TrueSkillService, deps.Templates)

	// Main app TrueSkill routes - protected by unified Discord OAuth
	mux.HandleFunc("/trueskill/update-all", deps.Auth.RequireAuth(trueskillHandler.UpdateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/update-user", deps.Auth.RequireAuth(trueskillHandler.UpdateUserTrueSkill))
	mux.HandleFunc("/trueskill/recalculate", deps.Auth.RequireAuth(trueskillHandler.RecalculateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/stats", deps.Auth.RequireAuth(trueskillHandler.GetTrueSkillStats))
}

func setupAPIRoutes(mux *http.ServeMux, deps *AppDependencies) {
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

func setupUSLRoutes(mux *http.ServeMux, deps *AppDependencies) {
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
	mux.HandleFunc("/usl/users/new", deps.Auth.RequireAuth(uslHandler.NewUserForm))
	mux.HandleFunc("/usl/users/create", deps.Auth.RequireAuth(uslHandler.CreateUser))
	mux.HandleFunc("/usl/users/edit", deps.Auth.RequireAuth(uslHandler.EditUserForm))
	mux.HandleFunc("/usl/users/update", deps.Auth.RequireAuth(uslHandler.UpdateUser))
	mux.HandleFunc("/usl/users/delete", deps.Auth.RequireAuth(uslHandler.DeleteUser))

	// USL Tracker Management Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/trackers", deps.Auth.RequireAuth(uslHandler.ListTrackers))
	mux.HandleFunc("/usl/trackers/new", deps.Auth.RequireAuth(uslHandler.NewTrackerForm))
	mux.HandleFunc("/usl/trackers/create", deps.Auth.RequireAuth(uslHandler.CreateTracker))

	// USL Admin Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/admin", deps.Auth.RequireAuth(uslHandler.AdminDashboard))
	mux.HandleFunc("/usl/import", deps.Auth.RequireAuth(uslHandler.ImportData))

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
