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

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Supabase client
	supabaseClient, err := supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.ServiceRoleKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(supabaseClient, cfg)
	trackerRepo := repositories.NewTrackerRepository(supabaseClient, cfg)

	// Initialize services
	percentileConverter := services.NewPercentileConverter(cfg)
	mmrCalculator := services.NewMMRCalculator(cfg, percentileConverter)
	enhancedUncertainty := services.NewEnhancedUncertaintyCalculator(cfg)
	dataTransformation := services.NewDataTransformationService()

	trueSkillService := services.NewUserTrueSkillService(
		trackerRepo,
		userRepo,
		mmrCalculator,
		enhancedUncertainty,
		dataTransformation,
		cfg,
	)

	// Parse templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, templates)
	trackerHandler := handlers.NewTrackerHandler(trackerRepo, templates)
	trueskillHandler := handlers.NewTrueSkillHandler(trueSkillService, templates)

	// Setup routes
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	})

	// User management routes
	mux.HandleFunc("/users", userHandler.ListUsers)
	mux.HandleFunc("/users/new", userHandler.NewUserForm)
	mux.HandleFunc("/users/create", userHandler.CreateUser)
	mux.HandleFunc("/users/edit", userHandler.EditUserForm)
	mux.HandleFunc("/users/update", userHandler.UpdateUser)
	mux.HandleFunc("/users/delete", userHandler.DeleteUser)
	mux.HandleFunc("/users/search", userHandler.SearchUsers)

	// Tracker management routes
	mux.HandleFunc("/trackers", trackerHandler.ListTrackers)
	mux.HandleFunc("/trackers/new", trackerHandler.NewTrackerForm)
	mux.HandleFunc("/trackers/create", trackerHandler.CreateTracker)
	mux.HandleFunc("/trackers/edit", trackerHandler.EditTrackerForm)
	mux.HandleFunc("/trackers/update", trackerHandler.UpdateTracker)
	mux.HandleFunc("/trackers/delete", trackerHandler.DeleteTracker)
	mux.HandleFunc("/trackers/search", trackerHandler.SearchTrackers)

	// TrueSkill calculation routes
	mux.HandleFunc("/trueskill/update-all", trueskillHandler.UpdateAllUserTrueSkill)
	mux.HandleFunc("/trueskill/update-user", trueskillHandler.UpdateUserTrueSkill)
	mux.HandleFunc("/trueskill/recalculate", trueskillHandler.RecalculateAllUserTrueSkill)
	mux.HandleFunc("/trueskill/stats", trueskillHandler.GetTrueSkillStats)

	// API routes (JSON responses)
	mux.HandleFunc("/api/users", userHandler.ListUsersAPI)
	mux.HandleFunc("/api/trackers", trackerHandler.ListTrackersAPI)
	mux.HandleFunc("/api/trueskill/update-all", trueskillHandler.UpdateAllUserTrueSkillAPI)

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	log.Printf("Supabase URL: %s", cfg.Supabase.URL)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
