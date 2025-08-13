package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// Config holds all configuration for the application
// Matches the configuration structure from the Google Apps Script project
type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	Supabase  SupabaseConfig  `json:"supabase"`
	TrueSkill TrueSkillConfig `json:"trueskill"`
	MMR       MMRConfig       `json:"mmr"`
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

// Load initializes configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Supabase: SupabaseConfig{
			URL:            getEnv("SUPABASE_URL", ""),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		},
		TrueSkill: TrueSkillConfig{
			InitialMu:            getEnvFloat("TRUESKILL_INITIAL_MU", 1500.0),
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

// GetTrueSkillSigmaRange returns the sigma min and max values
func (c *Config) GetTrueSkillSigmaRange() (float64, float64) {
	return c.TrueSkill.SigmaMax, c.TrueSkill.SigmaMin
}

// GetMMRConfig returns the MMR configuration
func (c *Config) GetMMRConfig() MMRConfig {
	return c.MMR
}
