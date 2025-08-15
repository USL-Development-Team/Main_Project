package usl

import (
	"fmt"
	"log/slog"
	"usl-server/internal/config"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// USLRepository provides simple, direct access to USL-specific tables
// No guild complexity, no abstractions - just basic CRUD operations
type USLRepository struct {
	client *supabase.Client
	config *config.Config
	logger *slog.Logger
}

// NewUSLRepository creates a new USL-specific repository
func NewUSLRepository(client *supabase.Client, config *config.Config, logger *slog.Logger) *USLRepository {
	return &USLRepository{
		client: client,
		config: config,
		logger: logger.With("component", "usl_repository"),
	}
}

func (r *USLRepository) GetAllUsers() ([]*USLUser, error) {
	var users []*USLUser

	_, err := r.client.From("usl_users").
		Select("*", "", false).
		Order("trueskill_mu", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func (r *USLRepository) SearchUsers(query string) ([]*USLUser, error) {
	var users []*USLUser

	_, err := r.client.From("usl_users").
		Select("*", "", false).
		Or(fmt.Sprintf("name.ilike.%%%s%%,discord_id.eq.%s", query, query), "").
		Order("trueskill_mu", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

func (r *USLRepository) GetUserByID(id int64) (*USLUser, error) {
	var user USLUser

	_, err := r.client.From("usl_users").
		Select("*", "", false).
		Eq("id", fmt.Sprintf("%d", id)).
		Single().
		ExecuteTo(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}

	return &user, nil
}

func (r *USLRepository) GetUserByDiscordID(discordID string) (*USLUser, error) {
	var user USLUser

	_, err := r.client.From("usl_users").
		Select("*", "", false).
		Eq("discord_id", discordID).
		Single().
		ExecuteTo(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by Discord ID %s: %w", discordID, err)
	}

	return &user, nil
}

func (r *USLRepository) CreateUser(name, discordID string, active, banned bool) (*USLUser, error) {
	insertData := map[string]interface{}{
		"name":       name,
		"discord_id": discordID,
		"active":     active,
		"banned":     banned,
	}

	var user USLUser
	_, err := r.client.From("usl_users").
		Insert(insertData, false, "", "", "").
		Single().
		ExecuteTo(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *USLRepository) UpdateUser(id int64, name string, active, banned bool) (*USLUser, error) {
	updateData := map[string]interface{}{
		"name":   name,
		"active": active,
		"banned": banned,
	}

	var user USLUser
	_, err := r.client.From("usl_users").
		Update(updateData, "", "").
		Eq("id", fmt.Sprintf("%d", id)).
		Single().
		ExecuteTo(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to update user %d: %w", id, err)
	}

	return &user, nil
}

func (r *USLRepository) DeleteUser(id int64) error {
	_, _, err := r.client.From("usl_users").
		Delete("", "").
		Eq("id", fmt.Sprintf("%d", id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete user %d: %w", id, err)
	}

	return nil
}

func (r *USLRepository) GetAllTrackers() ([]*USLUserTracker, error) {
	var trackers []*USLUserTracker

	_, err := r.client.From("usl_user_trackers").
		Select("*", "", false).
		Order("mmr", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&trackers)

	if err != nil {
		return nil, fmt.Errorf("failed to get all trackers: %w", err)
	}

	return trackers, nil
}

func (r *USLRepository) GetValidTrackers() ([]*USLUserTracker, error) {
	var trackers []*USLUserTracker

	_, err := r.client.From("usl_user_trackers").
		Select("*", "", false).
		Eq("valid", "true").
		Order("mmr", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&trackers)

	if err != nil {
		return nil, fmt.Errorf("failed to get valid trackers: %w", err)
	}

	return trackers, nil
}

func (r *USLRepository) GetTrackersByDiscordID(discordID string) ([]*USLUserTracker, error) {
	var trackers []*USLUserTracker

	_, err := r.client.From("usl_user_trackers").
		Select("*", "", false).
		Eq("discord_id", discordID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&trackers)

	if err != nil {
		return nil, fmt.Errorf("failed to get trackers for Discord ID %s: %w", discordID, err)
	}

	return trackers, nil
}

func (r *USLRepository) CreateTracker(tracker *USLUserTracker) (*USLUserTracker, error) {
	insertData := map[string]interface{}{
		"discord_id":                          tracker.DiscordID,
		"url":                                 tracker.URL,
		"ones_current_season_peak":            tracker.OnesCurrentSeasonPeak,
		"ones_previous_season_peak":           tracker.OnesPreviousSeasonPeak,
		"ones_all_time_peak":                  tracker.OnesAllTimePeak,
		"ones_current_season_games_played":    tracker.OnesCurrentSeasonGamesPlayed,
		"ones_previous_season_games_played":   tracker.OnesPreviousSeasonGamesPlayed,
		"twos_current_season_peak":            tracker.TwosCurrentSeasonPeak,
		"twos_previous_season_peak":           tracker.TwosPreviousSeasonPeak,
		"twos_all_time_peak":                  tracker.TwosAllTimePeak,
		"twos_current_season_games_played":    tracker.TwosCurrentSeasonGamesPlayed,
		"twos_previous_season_games_played":   tracker.TwosPreviousSeasonGamesPlayed,
		"threes_current_season_peak":          tracker.ThreesCurrentSeasonPeak,
		"threes_previous_season_peak":         tracker.ThreesPreviousSeasonPeak,
		"threes_all_time_peak":                tracker.ThreesAllTimePeak,
		"threes_current_season_games_played":  tracker.ThreesCurrentSeasonGamesPlayed,
		"threes_previous_season_games_played": tracker.ThreesPreviousSeasonGamesPlayed,
		"last_updated":                        tracker.LastUpdated,
		"valid":                               tracker.Valid,
		"mmr":                                 tracker.MMR,
	}

	var created USLUserTracker
	_, err := r.client.From("usl_user_trackers").
		Insert(insertData, false, "", "", "").
		Single().
		ExecuteTo(&created)

	if err != nil {
		return nil, fmt.Errorf("failed to create tracker: %w", err)
	}

	return &created, nil
}

func (r *USLRepository) GetUsersWithTrackers() ([]*USLUser, map[string][]*USLUserTracker, error) {
	users, err := r.GetAllUsers()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get users: %w", err)
	}

	allTrackers, err := r.GetAllTrackers()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get trackers: %w", err)
	}

	trackersByDiscord := make(map[string][]*USLUserTracker)
	for _, tracker := range allTrackers {
		trackersByDiscord[tracker.DiscordID] = append(trackersByDiscord[tracker.DiscordID], tracker)
	}

	return users, trackersByDiscord, nil
}

func (r *USLRepository) GetLeaderboard() ([]*USLUser, error) {
	var users []*USLUser

	_, err := r.client.From("usl_users").
		Select("*", "", false).
		Eq("active", "true").
		Order("trueskill_mu", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return users, nil
}

func (r *USLRepository) GetStats() (map[string]interface{}, error) {
	r.logger.Info("Getting USL stats")

	users, err := r.GetAllUsers()
	if err != nil {
		r.logger.Error("Failed to get users for stats", "error", err)
		return nil, fmt.Errorf("failed to get users for stats: %w", err)
	}
	r.logger.Info("Retrieved users for stats", "count", len(users))

	trackers, err := r.GetAllTrackers()
	if err != nil {
		r.logger.Error("Failed to get trackers for stats", "error", err)
		return nil, fmt.Errorf("failed to get trackers for stats: %w", err)
	}
	r.logger.Info("Retrieved trackers for stats", "count", len(trackers))

	activeUsers := 0
	bannedUsers := 0
	validTrackers := 0
	totalTrueSkillMu := 0.0

	for _, user := range users {
		if user.Active {
			activeUsers++
		}
		if user.Banned {
			bannedUsers++
		}
		totalTrueSkillMu += user.TrueSkillMu
	}

	for _, tracker := range trackers {
		if tracker.Valid {
			validTrackers++
		}
	}

	averageTrueSkillMu := 0.0
	if len(users) > 0 {
		averageTrueSkillMu = totalTrueSkillMu / float64(len(users))
	}

	return map[string]interface{}{
		"total_users":          len(users),
		"active_users":         activeUsers,
		"banned_users":         bannedUsers,
		"total_trackers":       len(trackers),
		"valid_trackers":       validTrackers,
		"average_trueskill_mu": averageTrueSkillMu,
	}, nil
}
