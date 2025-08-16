package services

import (
	"fmt"
	"log"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// UserTrueSkillService manages TrueSkill calculations and updates for individual users and batch operations.
// Exact port of JavaScript UserTrueSkillService with dependency injection pattern.
//
// Service Responsibilities:
// - Individual user TrueSkill calculation from tracker data
// - Batch TrueSkill updates for all users
// - Default TrueSkill assignment for users without trackers
// - TrueSkill recalculation workflows
type UserTrueSkillService struct {
	// Service dependencies - injected at runtime
	trackerRepo               *repositories.TrackerRepository
	userRepo                  *repositories.UserRepository
	percentileCalculator      *MMRCalculator
	enhancedUncertainty       *EnhancedUncertaintyCalculator
	dataTransformationService *DataTransformationService
	config                    *config.Config
}

// BatchUpdateResult represents batch processing results
type BatchUpdateResult struct {
	ProcessedCount    int                        `json:"processedCount"`
	TrackerBasedCount int                        `json:"trackerBasedCount"`
	DefaultCount      int                        `json:"defaultCount"`
	Errors            []TrueSkillProcessingError `json:"errors"`
}

// TrueSkillProcessingError represents errors during batch processing
type TrueSkillProcessingError struct {
	User      string `json:"user"`
	DiscordID string `json:"discordId"`
	Error     string `json:"error"`
}

// TrueSkillUpdateResult represents individual user update results
type TrueSkillUpdateResult struct {
	Success         bool                  `json:"success"`
	HadTrackers     bool                  `json:"hadTrackers"`
	TrueSkillResult *TrueSkillCalculation `json:"trueSkillResult,omitempty"`
	Error           string                `json:"error,omitempty"`
}

// TrueSkillCalculation represents TrueSkill calculation results
type TrueSkillCalculation struct {
	Mu          float64                `json:"mu"`
	Sigma       float64                `json:"sigma"`
	SkillResult *PercentileSkillResult `json:"skillResult"`
	LastUpdated time.Time              `json:"lastUpdated"`
}

// NewUserTrueSkillService creates a new user TrueSkill service instance
func NewUserTrueSkillService(
	trackerRepo *repositories.TrackerRepository,
	userRepo *repositories.UserRepository,
	percentileCalculator *MMRCalculator,
	enhancedUncertainty *EnhancedUncertaintyCalculator,
	dataTransformationService *DataTransformationService,
	config *config.Config,
) *UserTrueSkillService {
	return &UserTrueSkillService{
		trackerRepo:               trackerRepo,
		userRepo:                  userRepo,
		percentileCalculator:      percentileCalculator,
		enhancedUncertainty:       enhancedUncertainty,
		dataTransformationService: dataTransformationService,
		config:                    config,
	}
}

// UpdateAllUserTrueSkill updates TrueSkill values for all users
// Exact port of JavaScript updateAllUserTrueSkill() function
func (s *UserTrueSkillService) UpdateAllUserTrueSkill() (*BatchUpdateResult, error) {
	// Get all users from repository
	allUsers, err := s.userRepo.GetAllUsers(false) // Get all users, not just active
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	if len(allUsers) == 0 {
		log.Println("UserTrueSkillService: No users found in repository")
		return &BatchUpdateResult{
			ProcessedCount:    0,
			TrackerBasedCount: 0,
			DefaultCount:      0,
			Errors:            []TrueSkillProcessingError{},
		}, nil
	}

	log.Printf("UserTrueSkillService: Processing %d users for batch TrueSkill update", len(allUsers))

	var trackerBasedCount int
	var defaultCount int
	var errors []TrueSkillProcessingError
	var usersToUpdate []models.User

	// Calculate TrueSkill for all users and build batch update data
	for _, user := range allUsers {
		updatedUser := *user // Copy user data

		// Try to get trackers and calculate TrueSkill
		trackers, err := s.getUserTrackersForTrueSkill(user.DiscordID)
		if err != nil || len(trackers) == 0 {
			// No trackers - use existing values or defaults
			defaultCount++
			if !user.HasTrueSkillData() {
				// User has no TrueSkill data, set defaults
				s.setDefaultTrueSkillValues(&updatedUser)
			}
			updatedUser.TrueSkillLastUpdated = time.Now()
			log.Printf("Keeping existing TrueSkill for %s: μ=%.1f", user.Name, updatedUser.TrueSkillMu)
		} else {
			// User has trackers - calculate new TrueSkill values
			trackerData, err := s.dataTransformationService.PrepareTrackerDataForCalculation(trackers[0])
			if err != nil {
				log.Printf("Failed to prepare tracker data for %s: %v", user.Name, err)
				errors = append(errors, TrueSkillProcessingError{
					User:      user.Name,
					DiscordID: user.DiscordID,
					Error:     fmt.Sprintf("failed to prepare tracker data: %v", err),
				})
				// Fall back to existing values
				defaultCount++
				updatedUser.TrueSkillLastUpdated = time.Now()
				continue
			}

			trueSkillResult, err := s.calculateTrueSkillValues(trackerData)
			if err != nil {
				log.Printf("Failed to calculate TrueSkill for %s: %v", user.Name, err)
				errors = append(errors, TrueSkillProcessingError{
					User:      user.Name,
					DiscordID: user.DiscordID,
					Error:     fmt.Sprintf("failed to calculate TrueSkill: %v", err),
				})
				// Fall back to existing values
				defaultCount++
				updatedUser.TrueSkillLastUpdated = time.Now()
				continue
			}

			trackerBasedCount++
			updatedUser.TrueSkillMu = trueSkillResult.Mu
			updatedUser.TrueSkillSigma = trueSkillResult.Sigma
			updatedUser.TrueSkillLastUpdated = trueSkillResult.LastUpdated
			log.Printf("Calculated TrueSkill for %s: μ=%.1f", user.Name, trueSkillResult.Mu)
		}

		usersToUpdate = append(usersToUpdate, updatedUser)
	}

	// Batch update all users
	log.Printf("UserTrueSkillService: Batch updating %d users with new TrueSkill values", len(usersToUpdate))
	successCount, err := s.userRepo.BatchUpdateTrueSkill(usersToUpdate)
	if err != nil {
		return nil, fmt.Errorf("batch update failed: %w", err)
	}

	log.Printf("UserTrueSkillService: Batch update completed successfully, updated %d users", successCount)

	result := &BatchUpdateResult{
		ProcessedCount:    len(allUsers),
		TrackerBasedCount: trackerBasedCount,
		DefaultCount:      defaultCount,
		Errors:            errors,
	}

	s.reportBatchUpdateResults(result)
	return result, nil
}

// UpdateUserTrueSkillFromTrackerData updates TrueSkill for a single user from provided tracker data
// This method bypasses the repository layer and accepts TrackerData directly
func (s *UserTrueSkillService) UpdateUserTrueSkillFromTrackerData(trackerData *TrackerData) *TrueSkillUpdateResult {
	// Validate tracker data
	if err := s.dataTransformationService.ValidateTrackerData(trackerData); err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("invalid tracker data: %v", err),
		}
	}

	// Calculate TrueSkill values directly from provided tracker data
	trueSkillResult, err := s.calculateTrueSkillValues(trackerData)
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to calculate TrueSkill: %v", err),
		}
	}

	// Update user with new TrueSkill values
	err = s.updateUserWithTrueSkillValues(trackerData.DiscordID, trueSkillResult)
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to update user: %v", err),
		}
	}

	log.Printf("UserTrueSkillService: Updated TrueSkill for user %s: μ=%.1f, σ=%.2f (from TrackerData)",
		trackerData.DiscordID, trueSkillResult.Mu, trueSkillResult.Sigma)

	return &TrueSkillUpdateResult{
		Success:         true,
		HadTrackers:     true,
		TrueSkillResult: trueSkillResult,
	}
}

// CalculateTrueSkillFromTrackerData calculates TrueSkill values from tracker data without database operations
// This is a pure calculation method that delegates to the internal calculation logic
func (s *UserTrueSkillService) CalculateTrueSkillFromTrackerData(trackerData *TrackerData) *TrueSkillUpdateResult {
	// Validate tracker data
	if err := s.dataTransformationService.ValidateTrackerData(trackerData); err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("invalid tracker data: %v", err),
		}
	}

	// Calculate TrueSkill values directly from provided tracker data
	trueSkillResult, err := s.calculateTrueSkillValues(trackerData)
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to calculate TrueSkill: %v", err),
		}
	}

	log.Printf("UserTrueSkillService: Calculated TrueSkill for user %s: μ=%.1f, σ=%.2f (calculation only)",
		trackerData.DiscordID, trueSkillResult.Mu, trueSkillResult.Sigma)

	return &TrueSkillUpdateResult{
		Success:         true,
		HadTrackers:     true,
		TrueSkillResult: trueSkillResult,
	}
}

// UpdateUserTrueSkillFromTrackers updates TrueSkill for a single user from their tracker data
// Exact port of JavaScript updateUserTrueSkillFromTrackers() function
func (s *UserTrueSkillService) UpdateUserTrueSkillFromTrackers(discordID string) *TrueSkillUpdateResult {
	trackers, err := s.getUserTrackersForTrueSkill(discordID)
	if err != nil || len(trackers) == 0 {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: false,
		}
	}

	trackerData, err := s.dataTransformationService.PrepareTrackerDataForCalculation(trackers[0])
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to prepare tracker data: %v", err),
		}
	}

	trueSkillResult, err := s.calculateTrueSkillValues(trackerData)
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to calculate TrueSkill: %v", err),
		}
	}

	// Update user with new TrueSkill values
	err = s.updateUserWithTrueSkillValues(discordID, trueSkillResult)
	if err != nil {
		return &TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to update user: %v", err),
		}
	}

	log.Printf("UserTrueSkillService: Updated TrueSkill for user %s: μ=%.1f, σ=%.2f (percentile-based)",
		discordID, trueSkillResult.Mu, trueSkillResult.Sigma)

	return &TrueSkillUpdateResult{
		Success:         true,
		HadTrackers:     true,
		TrueSkillResult: trueSkillResult,
	}
}

// RecalculateAllUserTrueSkill recalculates TrueSkill values for all users
// Exact port of JavaScript recalculateAllUserTrueSkill() function
func (s *UserTrueSkillService) RecalculateAllUserTrueSkill() (*BatchUpdateResult, error) {
	log.Println("UserTrueSkillService: Recalculating all user TrueSkill values (delegating to UpdateAllUserTrueSkill)")
	return s.UpdateAllUserTrueSkill()
}

// getUserTrackersForTrueSkill gets valid trackers for a user's TrueSkill calculation
// Exact port of JavaScript _getUserTrackersForTrueSkill() function
func (s *UserTrueSkillService) getUserTrackersForTrueSkill(discordID string) ([]*models.UserTracker, error) {
	trackers, err := s.trackerRepo.GetTrackersByDiscordID(discordID, true) // activeOnly = true
	if err != nil {
		log.Printf("UserTrueSkillService: Error getting trackers for user %s: %v", discordID, err)
		return nil, err
	}

	if len(trackers) == 0 {
		log.Printf("UserTrueSkillService: No trackers found for user %s", discordID)
		return nil, nil
	}

	return trackers, nil
}

// calculateTrueSkillValues calculates TrueSkill values from tracker data
// Exact port of JavaScript _calculateTrueSkillValues() function
func (s *UserTrueSkillService) calculateTrueSkillValues(trackerData *TrackerData) (*TrueSkillCalculation, error) {
	// Calculate percentile-based TrueSkill seeding using structured object format
	playerData := PlayerData{
		Ones: PlaylistData{
			Current:  PlaylistSeasonData{MMR: trackerData.OnesCurrentPeak, Games: trackerData.OnesCurrentGames},
			Previous: PlaylistSeasonData{MMR: trackerData.OnesPreviousPeak, Games: trackerData.OnesPreviousGames},
		},
		Twos: PlaylistData{
			Current:  PlaylistSeasonData{MMR: trackerData.TwosCurrentPeak, Games: trackerData.TwosCurrentGames},
			Previous: PlaylistSeasonData{MMR: trackerData.TwosPreviousPeak, Games: trackerData.TwosPreviousGames},
		},
		Threes: PlaylistData{
			Current:  PlaylistSeasonData{MMR: trackerData.ThreesCurrentPeak, Games: trackerData.ThreesCurrentGames},
			Previous: PlaylistSeasonData{MMR: trackerData.ThreesPreviousPeak, Games: trackerData.ThreesPreviousGames},
		},
	}

	skillResult := s.percentileCalculator.CalculatePercentileBasedSkill(playerData)

	trueskillSigma, err := s.enhancedUncertainty.CalculateEnhancedUncertainty(trackerData)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate enhanced uncertainty: %w", err)
	}

	return &TrueSkillCalculation{
		Mu:          skillResult.TrueskillMu,
		Sigma:       trueskillSigma,
		SkillResult: &skillResult,
		LastUpdated: time.Now(),
	}, nil
}

// updateUserWithTrueSkillValues updates a user's TrueSkill values in the database
// Exact port of JavaScript _updateUserWithTrueSkillValues() function
func (s *UserTrueSkillService) updateUserWithTrueSkillValues(discordID string, trueSkillResult *TrueSkillCalculation) error {
	return s.userRepo.UpdateUserTrueSkill(
		discordID,
		trueSkillResult.Mu,
		trueSkillResult.Sigma,
		&trueSkillResult.LastUpdated,
	)
}

// setDefaultTrueSkillValues sets default TrueSkill values for users without tracker data
// Exact port of JavaScript _setDefaultTrueSkillValues() function
func (s *UserTrueSkillService) setDefaultTrueSkillValues(user *models.User) {
	// For users without tracker data, use reasonable defaults
	defaultMu, defaultSigma := s.config.GetTrueSkillDefaults()

	user.TrueSkillMu = defaultMu
	user.TrueSkillSigma = defaultSigma

	log.Printf("UserTrueSkillService: Set default TrueSkill for user %s: μ=%.1f, σ=%.1f",
		user.DiscordID, defaultMu, defaultSigma)
}

// reportBatchUpdateResults reports the results of batch processing
// Exact port of JavaScript _reportBatchUpdateResults() function
func (s *UserTrueSkillService) reportBatchUpdateResults(result *BatchUpdateResult) {
	log.Printf("UserTrueSkillService: Successfully processed %d users (%d given default values)",
		result.ProcessedCount, result.DefaultCount)

	if len(result.Errors) > 0 {
		log.Printf("UserTrueSkillService: %d errors occurred during batch processing", len(result.Errors))
		for _, err := range result.Errors {
			log.Printf("  - User %s (%s): %s", err.User, err.DiscordID, err.Error)
		}
	}

	log.Printf("TrueSkill update complete! Processed %d users. %d calculated from trackers, %d given default values.",
		result.ProcessedCount, result.TrackerBasedCount, result.DefaultCount)
}

// GetTrueSkillStats returns service statistics
// Exact port of JavaScript getTrueSkillStats() function
func (s *UserTrueSkillService) GetTrueSkillStats() map[string]interface{} {
	return map[string]interface{}{
		"serviceName": "UserTrueSkillService",
		"version":     "1.0.0",
		"dependencies": map[string]bool{
			"trackerRepo":               s.trackerRepo != nil,
			"userRepo":                  s.userRepo != nil,
			"percentileCalculator":      s.percentileCalculator != nil,
			"enhancedUncertainty":       s.enhancedUncertainty != nil,
			"dataTransformationService": s.dataTransformationService != nil,
		},
	}
}
