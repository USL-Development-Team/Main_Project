package main

import (
	"html/template"
	"log"
	"net/http"
	"usl-server/internal/config"
	"usl-server/internal/handlers"
	"usl-server/internal/repositories"
	"usl-server/internal/services"

	"github.com/supabase-community/supabase-go"
)

type AppDependencies struct {
	Config           *config.Config
	UserRepo         *repositories.UserRepository
	TrackerRepo      *repositories.TrackerRepository
	TrueSkillService *services.UserTrueSkillService
	Templates        *template.Template
}

func main() {
	dependencies := initializeApplication()
	server := setupHTTPServer(dependencies)
	startServer(server, dependencies.Config)
}

func initializeApplication() *AppDependencies {
	configuration := loadConfiguration()
	supabaseClient := createSupabaseClient(configuration)
	repositories := setupRepositories(supabaseClient, configuration)
	services := setupServices(configuration, repositories)
	templates := loadTemplates()

	return &AppDependencies{
		Config:           configuration,
		UserRepo:         repositories.UserRepo,
		TrackerRepo:      repositories.TrackerRepo,
		TrueSkillService: services,
		Templates:        templates,
	}
}

func loadConfiguration() *config.Config {
	configuration, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return configuration
}

func createSupabaseClient(configuration *config.Config) *supabase.Client {
	client, err := supabase.NewClient(
		configuration.Supabase.URL,
		configuration.Supabase.ServiceRoleKey,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}
	return client
}

type RepositoryCollection struct {
	UserRepo    *repositories.UserRepository
	TrackerRepo *repositories.TrackerRepository
}

func setupRepositories(client *supabase.Client, config *config.Config) *RepositoryCollection {
	return &RepositoryCollection{
		UserRepo:    repositories.NewUserRepository(client, config),
		TrackerRepo: repositories.NewTrackerRepository(client, config),
	}
}

func setupServices(config *config.Config, repos *RepositoryCollection) *services.UserTrueSkillService {
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

func loadTemplates() *template.Template {
	return template.Must(template.ParseGlob("templates/*.html"))
}

func setupHTTPServer(deps *AppDependencies) *http.ServeMux {
	mux := http.NewServeMux()

	setupStaticRoutes(mux)
	setupHomeRoute(mux)
	setupUserRoutes(mux, deps)
	setupTrackerRoutes(mux, deps)
	setupTrueSkillRoutes(mux, deps)
	setupAPIRoutes(mux, deps)

	return mux
}

func setupStaticRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}

func setupHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	})
}

func setupUserRoutes(mux *http.ServeMux, deps *AppDependencies) {
	userHandler := handlers.NewUserHandler(deps.UserRepo, deps.Templates)

	mux.HandleFunc("/users", userHandler.ListUsers)
	mux.HandleFunc("/users/new", userHandler.NewUserForm)
	mux.HandleFunc("/users/create", userHandler.CreateUser)
	mux.HandleFunc("/users/edit", userHandler.EditUserForm)
	mux.HandleFunc("/users/update", userHandler.UpdateUser)
	mux.HandleFunc("/users/delete", userHandler.DeleteUser)
	mux.HandleFunc("/users/search", userHandler.SearchUsers)
}

func setupTrackerRoutes(mux *http.ServeMux, deps *AppDependencies) {
	trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.Templates)

	mux.HandleFunc("/trackers", trackerHandler.ListTrackers)
	mux.HandleFunc("/trackers/new", trackerHandler.NewTrackerForm)
	mux.HandleFunc("/trackers/create", trackerHandler.CreateTracker)
	mux.HandleFunc("/trackers/edit", trackerHandler.EditTrackerForm)
	mux.HandleFunc("/trackers/update", trackerHandler.UpdateTracker)
	mux.HandleFunc("/trackers/delete", trackerHandler.DeleteTracker)
	mux.HandleFunc("/trackers/search", trackerHandler.SearchTrackers)
}

func setupTrueSkillRoutes(mux *http.ServeMux, deps *AppDependencies) {
	trueskillHandler := handlers.NewTrueSkillHandler(deps.TrueSkillService, deps.Templates)

	mux.HandleFunc("/trueskill/update-all", trueskillHandler.UpdateAllUserTrueSkill)
	mux.HandleFunc("/trueskill/update-user", trueskillHandler.UpdateUserTrueSkill)
	mux.HandleFunc("/trueskill/recalculate", trueskillHandler.RecalculateAllUserTrueSkill)
	mux.HandleFunc("/trueskill/stats", trueskillHandler.GetTrueSkillStats)
}

func setupAPIRoutes(mux *http.ServeMux, deps *AppDependencies) {
	userHandler := handlers.NewUserHandler(deps.UserRepo, deps.Templates)
	trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.Templates)
	trueskillHandler := handlers.NewTrueSkillHandler(deps.TrueSkillService, deps.Templates)

	mux.HandleFunc("/api/users", userHandler.ListUsersAPI)
	mux.HandleFunc("/api/trackers", trackerHandler.ListTrackersAPI)
	mux.HandleFunc("/api/trueskill/update-all", trueskillHandler.UpdateAllUserTrueSkillAPI)
}

func startServer(server *http.ServeMux, config *config.Config) {
	serverAddress := config.Server.Host + ":" + config.Server.Port

	log.Printf("Starting USL server on %s", serverAddress)
	log.Printf("Supabase URL: %s", config.Supabase.URL)
	log.Printf("Visit http://%s to access the application", serverAddress)

	if err := http.ListenAndServe(serverAddress, server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
