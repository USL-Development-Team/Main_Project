package main

import (
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
	"usl-server/internal/templates"
	usl "usl-server/internal/usl"
	uslHandlers "usl-server/internal/usl/handlers"

	"github.com/supabase-community/supabase-go"
)

const (
	templatePattern = "templates/*.html"
	healthResponse  = `{"status":"healthy","service":"usl-server"}`
)

type ApplicationContext struct {
	Config *config.Config
	Logger *slog.Logger

	Auth *auth.DiscordAuth

	UserRepo    *repositories.UserRepository
	TrackerRepo *repositories.TrackerRepository
	GuildRepo   *repositories.GuildRepository

	TrueSkillService *services.UserTrueSkillService

	Templates *template.Template
}

func main() {
	logger := logger.SetupLogger(slog.LevelDebug, "logs/server.log")
	logger.Info("Starting USL server application")

	app := initializeApplication(logger)
	server := setupHTTPServer(app)
	startServer(server, app.Config, app.Logger)
}

func initializeApplication(logger *slog.Logger) *ApplicationContext {
	appConfig := loadConfiguration(logger)
	supabaseClient := createSupabaseClient(appConfig, logger)
	repositories := setupRepositories(supabaseClient, appConfig, logger)
	services := setupServices(appConfig, repositories, logger)
	templates := loadTemplates(logger)

	envConfig := config.GetEnvironmentConfig()

	discordAuth := auth.NewDiscordAuth(supabaseClient, appConfig.USL.AdminDiscordIDs,
		appConfig.Supabase.URL, appConfig.Supabase.PublicURL, appConfig.Supabase.AnonKey, envConfig)

	return &ApplicationContext{
		Config:           appConfig,
		UserRepo:         repositories.UserRepo,
		TrackerRepo:      repositories.TrackerRepo,
		GuildRepo:        repositories.GuildRepo,
		TrueSkillService: services,
		Templates:        templates,
		Logger:           logger,
		Auth:             discordAuth,
	}
}

func loadConfiguration(logger *slog.Logger) *config.Config {
	logger.Info("Loading configuration")
	appConfig, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully")
	return appConfig
}

func createSupabaseClient(appConfig *config.Config, logger *slog.Logger) *supabase.Client {
	logger.Info("Initializing Supabase client", "url", appConfig.Supabase.URL)
	client, err := supabase.NewClient(
		appConfig.Supabase.URL,
		appConfig.Supabase.ServiceRoleKey,
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

func setupRepositories(client *supabase.Client, appConfig *config.Config, logger *slog.Logger) *RepositoryCollection {
	logger.Info("Setting up repositories")
	return &RepositoryCollection{
		UserRepo:    repositories.NewUserRepository(client, appConfig),
		TrackerRepo: repositories.NewTrackerRepository(client, appConfig),
		GuildRepo:   repositories.NewGuildRepository(client, appConfig),
	}
}

func setupServices(appConfig *config.Config, repos *RepositoryCollection, logger *slog.Logger) *services.UserTrueSkillService {
	logger.Info("Setting up services")
	percentileConverter := services.NewPercentileConverter(appConfig)
	mmrCalculator := services.NewMMRCalculator(appConfig, percentileConverter)
	uncertaintyCalculator := services.NewEnhancedUncertaintyCalculator(appConfig)
	dataTransformationService := services.NewDataTransformationService()

	return services.NewUserTrueSkillService(
		repos.TrackerRepo,
		repos.UserRepo,
		mmrCalculator,
		uncertaintyCalculator,
		dataTransformationService,
		appConfig,
	)
}

func createTemplateFunctions() template.FuncMap {
	return templates.TemplateFunctions()
}

func loadTemplates(logger *slog.Logger) *template.Template {
	logger.Info("Loading templates")

	templateEngine := template.New("app")
	templateEngine = templateEngine.Funcs(createTemplateFunctions())

	matches, err := filepath.Glob(templatePattern)
	if err != nil {
		logger.Error("Failed to glob templates", "pattern", templatePattern, "error", err)
		os.Exit(1)
	}

	if len(matches) > 0 {
		logger.Info("Loading templates", "pattern", templatePattern, "count", len(matches))
		templateEngine = template.Must(templateEngine.ParseFiles(matches...))
	}

	logger.Info("Templates loaded successfully")
	return templateEngine
}

func setupHTTPServer(app *ApplicationContext) http.Handler {
	mux := http.NewServeMux()

	setupStaticRoutes(mux)
	setupHealthRoute(mux)
	setupAuthRoutes(mux, app)
	setupHomeRoute(mux)
	setupGuildAwareRoutes(mux, app)
	setupUserRoutes(mux, app)
	setupTrackerRoutes(mux, app)
	setupTrueSkillRoutes(mux, app)
	setupAPIRoutes(mux, app)
	setupUSLRoutes(mux, app)

	guildMiddleware := middleware.NewGuildContextMiddleware(app.GuildRepo, app.Logger)
	handler := middleware.LoggingMiddleware(app.Logger)(mux)
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
		if _, err := w.Write([]byte(healthResponse)); err != nil {
			log.Printf("Failed to write health response: %v", err)
		}
	})
}

func setupGuildAwareRoutes(mux *http.ServeMux, app *ApplicationContext) {
	// TODO: Add support for dynamic guild slugs when needed
}

func setupAuthRoutes(mux *http.ServeMux, app *ApplicationContext) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
	})
	mux.HandleFunc("/usl/login", app.Auth.LoginForm)
	mux.HandleFunc("/auth/callback", app.Auth.AuthCallback)
	mux.HandleFunc("/auth/process", app.Auth.ProcessTokens)
	mux.HandleFunc("/logout", app.Auth.Logout)
	mux.HandleFunc("/usl/logout", app.Auth.Logout)
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

	if _, err := w.Write([]byte(templates.NotFoundHTML)); err != nil {
		log.Printf("Failed to write 404 response: %v", err)
	}
}

func setupUserRoutes(mux *http.ServeMux, app *ApplicationContext) {
	mux.HandleFunc("/users", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/new", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/new", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/create", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/create", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/edit", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/edit", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/update", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/update", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/delete", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/delete", http.StatusMovedPermanently)
	}))
	mux.HandleFunc("/users/search", app.Auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/usl/users/search", http.StatusMovedPermanently)
	}))
}

func setupTrackerRoutes(mux *http.ServeMux, app *ApplicationContext) {
	trackerHandler := handlers.NewTrackerHandler(app.TrackerRepo, app.TrueSkillService, app.Templates)

	mux.HandleFunc("/trackers", app.Auth.RequireAuth(trackerHandler.ListTrackers))
	mux.HandleFunc("/trackers/new", app.Auth.RequireAuth(trackerHandler.NewTrackerForm))
	mux.HandleFunc("/trackers/create", app.Auth.RequireAuth(trackerHandler.CreateTracker))
	mux.HandleFunc("/trackers/edit", app.Auth.RequireAuth(trackerHandler.EditTrackerForm))
	mux.HandleFunc("/trackers/update", app.Auth.RequireAuth(trackerHandler.UpdateTracker))
	mux.HandleFunc("/trackers/delete", app.Auth.RequireAuth(trackerHandler.DeleteTracker))
	mux.HandleFunc("/trackers/search", app.Auth.RequireAuth(trackerHandler.SearchTrackers))
}

func setupTrueSkillRoutes(mux *http.ServeMux, app *ApplicationContext) {
	trueskillHandler := handlers.NewTrueSkillHandler(app.TrueSkillService, app.Templates)

	mux.HandleFunc("/trueskill/update-all", app.Auth.RequireAuth(trueskillHandler.UpdateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/update-user", app.Auth.RequireAuth(trueskillHandler.UpdateUserTrueSkill))
	mux.HandleFunc("/trueskill/recalculate", app.Auth.RequireAuth(trueskillHandler.RecalculateAllUserTrueSkill))
	mux.HandleFunc("/trueskill/stats", app.Auth.RequireAuth(trueskillHandler.GetTrueSkillStats))
}

func setupAPIRoutes(mux *http.ServeMux, app *ApplicationContext) {
	userHandler := handlers.NewUserHandler(app.UserRepo, app.Templates)
	trackerHandler := handlers.NewTrackerHandler(app.TrackerRepo, app.TrueSkillService, app.Templates)
	trueskillHandler := handlers.NewTrueSkillHandler(app.TrueSkillService, app.Templates)

	v2UsersHandler := uslHandlers.NewV2UsersHandler(app.UserRepo)
	v2TrackersHandler := uslHandlers.NewV2TrackersHandler(app.TrackerRepo)

	mux.HandleFunc("/api/users", app.Auth.RequireAuth(userHandler.ListUsersAPI))
	mux.HandleFunc("/api/trackers", app.Auth.RequireAuth(trackerHandler.ListTrackersAPI))
	mux.HandleFunc("/api/trueskill/update-all", app.Auth.RequireAuth(trueskillHandler.UpdateAllUserTrueSkillAPI))

	mux.HandleFunc("/api/v2/users", app.Auth.RequireAuth(v2UsersHandler.HandleUsers))
	mux.HandleFunc("/api/v2/users/bulk", app.Auth.RequireAuth(v2UsersHandler.HandleUsersBulk))
	mux.HandleFunc("/api/v2/trackers", app.Auth.RequireAuth(v2TrackersHandler.HandleTrackers))
	mux.HandleFunc("/api/v2/trackers/bulk", app.Auth.RequireAuth(v2TrackersHandler.HandleTrackersBulk))
}

func setupUSLRoutes(mux *http.ServeMux, app *ApplicationContext) {
	// TEMPORARY: USL migration handler - will be deleted after migration
	// Create dedicated Supabase client for USL operations
	supabaseClient, err := supabase.NewClient(
		app.Config.Supabase.URL,
		app.Config.Supabase.ServiceRoleKey,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create USL Supabase client: %v", err)
	}

	uslRepo := usl.NewUSLRepository(supabaseClient, app.Config, app.Logger)
	// NOTE: USL handlers no longer need their own auth - they use the unified auth
	uslHandler := uslHandlers.NewMigrationHandler(uslRepo, app.Templates, app.TrueSkillService, app.Config)

	// USL Main Routes (redirect to login or admin based on auth status)
	mux.HandleFunc("/usl/", func(w http.ResponseWriter, r *http.Request) {
		if app.Auth.IsAuthenticated(r) {
			http.Redirect(w, r, "/usl/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
		}
	})

	// USL User Management Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/users", app.Auth.RequireAuth(uslHandler.ListUsers))
	mux.HandleFunc("/usl/users/search", app.Auth.RequireAuth(uslHandler.SearchUsers))
	mux.HandleFunc("/usl/users/detail", app.Auth.RequireAuth(uslHandler.UserDetail))
	mux.HandleFunc("/usl/users/new", app.Auth.RequireAuth(uslHandler.NewUserForm))
	mux.HandleFunc("/usl/users/create", app.Auth.RequireAuth(uslHandler.CreateUser))
	mux.HandleFunc("/usl/users/edit", app.Auth.RequireAuth(uslHandler.EditUserForm))
	mux.HandleFunc("/usl/users/update", app.Auth.RequireAuth(uslHandler.UpdateUser))
	mux.HandleFunc("/usl/users/delete", app.Auth.RequireAuth(uslHandler.DeleteUser))
	mux.HandleFunc("/usl/users/update-trueskill", app.Auth.RequireAuth(uslHandler.UpdateUserTrueSkill))

	// USL Tracker Management Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/trackers", app.Auth.RequireAuth(uslHandler.ListTrackers))
	mux.HandleFunc("/usl/trackers/search", app.Auth.RequireAuth(uslHandler.SearchTrackers))
	mux.HandleFunc("/usl/trackers/detail", app.Auth.RequireAuth(uslHandler.TrackerDetail))
	mux.HandleFunc("/usl/trackers/new", app.Auth.RequireAuth(uslHandler.NewTrackerForm))
	mux.HandleFunc("/usl/trackers/create", app.Auth.RequireAuth(uslHandler.CreateTracker))
	mux.HandleFunc("/usl/trackers/edit", app.Auth.RequireAuth(uslHandler.EditTrackerForm))
	mux.HandleFunc("/usl/trackers/update", app.Auth.RequireAuth(uslHandler.UpdateTracker))

	// USL Admin Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/admin", app.Auth.RequireAuth(uslHandler.AdminDashboard))

	// USL API Routes (protected by unified Discord OAuth)
	mux.HandleFunc("/usl/api/users", app.Auth.RequireAuth(uslHandler.ListUsersAPI))
	mux.HandleFunc("/usl/api/trackers", app.Auth.RequireAuth(uslHandler.ListTrackersAPI))
	mux.HandleFunc("/usl/api/leaderboard", app.Auth.RequireAuth(uslHandler.GetLeaderboardAPI))
}

func startServer(server http.Handler, appConfig *config.Config, logger *slog.Logger) {
	serverAddress := appConfig.Server.Host + ":" + appConfig.Server.Port

	logger.Info("Starting USL server",
		"address", serverAddress,
		"supabase_url", appConfig.Supabase.URL)

	if err := http.ListenAndServe(serverAddress, server); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
