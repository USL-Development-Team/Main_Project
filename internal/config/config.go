package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
// Matches the configuration structure from the Google Apps Script project
type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	Supabase  SupabaseConfig  `json:"supabase"`
	TrueSkill TrueSkillConfig `json:"trueskill"`
	MMR       MMRConfig       `json:"mmr"`
	USL       USLConfig       `json:"usl"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

type SupabaseConfig struct {
	URL            string `json:"url"`
	PublicURL      string `json:"public_url"` // For OAuth callbacks via ngrok
	AnonKey        string `json:"anon_key"`
	ServiceRoleKey string `json:"service_role_key"`
}

// TrueSkillConfig matches the configuration from TrueSkillConfig.js
type TrueSkillConfig struct {
	InitialMu            float64 `json:"initial_mu"`
	InitialSigma         float64 `json:"initial_sigma"`
	SigmaMin             float64 `json:"sigma_min"`
	SigmaMax             float64 `json:"sigma_max"`
	GamesForMaxCertainty int     `json:"games_for_max_certainty"`
}

// MMRConfig matches the configuration from MMRConfig.js
type MMRConfig struct {
	OnesWeight           float64 `json:"ones_weight"`
	TwosWeight           float64 `json:"twos_weight"`
	ThreesWeight         float64 `json:"threes_weight"`
	MinGamesThreshold    int     `json:"min_games_threshold"`
	CurrentSeasonWeight  float64 `json:"current_season_weight"`
	PreviousSeasonWeight float64 `json:"previous_season_weight"`
}

// USLConfig holds USL-specific configuration for temporary migration
type USLConfig struct {
	AdminDiscordIDs []string `json:"admin_discord_ids"`
}

// Load initializes configuration from environment variables
func Load() (*Config, error) {
	// Skip .env file loading if running on a platform that provides environment variables
	// Check for common platform environment indicators
	if !isPlatformEnvironment() {
		// Load .env file for local development
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}
	} else {
		log.Printf("Platform environment detected, skipping .env file loading")
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("PORT", getEnv("SERVER_PORT", "8080")), // Render provides PORT env var
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Supabase: SupabaseConfig{
			URL:            getEnv("SUPABASE_URL", ""),
			PublicURL:      getEnv("SUPABASE_PUBLIC_URL", getEnv("SUPABASE_URL", "")), // Fallback to URL if not set
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		},
		TrueSkill: TrueSkillConfig{
			InitialMu:            getEnvFloat("TRUESKILL_INITIAL_MU", 1000.0),
			InitialSigma:         getEnvFloat("TRUESKILL_INITIAL_SIGMA", 8.333),
			SigmaMin:             getEnvFloat("TRUESKILL_SIGMA_MIN", 2.5),
			SigmaMax:             getEnvFloat("TRUESKILL_SIGMA_MAX", 8.333),
			GamesForMaxCertainty: getEnvInt("TRUESKILL_GAMES_FOR_MAX_CERTAINTY", 1000),
		},
		MMR: MMRConfig{
			OnesWeight:           getEnvFloat("MMR_ONES_WEIGHT", 1.0),
			TwosWeight:           getEnvFloat("MMR_TWOS_WEIGHT", 1.5),
			ThreesWeight:         getEnvFloat("MMR_THREES_WEIGHT", 1.2),
			MinGamesThreshold:    getEnvInt("MMR_MIN_GAMES_THRESHOLD", 10),
			CurrentSeasonWeight:  getEnvFloat("MMR_CURRENT_SEASON_WEIGHT", 0.7),
			PreviousSeasonWeight: getEnvFloat("MMR_PREVIOUS_SEASON_WEIGHT", 0.3),
		},
		USL: USLConfig{
			AdminDiscordIDs: getEnvStringSlice("USL_ADMIN_DISCORD_IDS", []string{"679038415576104971", "354474826192388127"}),
		},
	}

	return config, nil
}

// Utility functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// GetTrueSkillDefaults returns the default TrueSkill values for new users
// Matches the default values from the Google Apps Script project
func (c *Config) GetTrueSkillDefaults() (float64, float64) {
	return c.TrueSkill.InitialMu, c.TrueSkill.InitialSigma
}

func (c *Config) GetTrueSkillSigmaRange() (float64, float64) {
	return c.TrueSkill.SigmaMax, c.TrueSkill.SigmaMin
}

func (c *Config) GetMMRConfig() MMRConfig {
	return c.MMR
}

// isPlatformEnvironment detects if running on a platform that provides environment variables
func isPlatformEnvironment() bool {
	// Check for common platform environment indicators
	platformVars := []string{
		"RENDER",             // Render.com
		"HEROKU",             // Heroku
		"VERCEL",             // Vercel
		"RAILWAY_PROJECT_ID", // Railway
		"FLY_APP_NAME",       // Fly.io
		"CF_INSTANCE_INDEX",  // Cloud Foundry
	}

	for _, envVar := range platformVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	// Also check if ENVIRONMENT is explicitly set to production
	// and key environment variables are already available
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		// Check if critical environment variables are already set
		requiredVars := []string{
			"SUPABASE_URL",
			"SUPABASE_ANON_KEY",
			"SUPABASE_SERVICE_ROLE_KEY",
		}

		allPresent := true
		for _, envVar := range requiredVars {
			if os.Getenv(envVar) == "" {
				allPresent = false
				break
			}
		}

		if allPresent {
			return true
		}
	}

	return false
}

// getEnvStringSlice parses a comma-separated string into a slice
func getEnvStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}
