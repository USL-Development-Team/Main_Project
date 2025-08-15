package handlers

// Shared API response messages for consistency across all handlers
const (
	// HTTP method errors
	msgMethodNotAllowed = "method not allowed"

	// Request validation errors
	msgInvalidRequestBody = "invalid request body"
	msgValidationFailed   = "validation failed"
	msgInvalidSortField   = "invalid sort field"

	// Operation errors
	msgBulkOperationFailed   = "bulk operation failed"
	msgFailedToGetUsers      = "failed to get users"
	msgFailedToCreateUser    = "failed to create user"
	msgUserAlreadyExists     = "user already exists"
	msgFailedToGetTrackers   = "failed to get trackers"
	msgFailedToCreateTracker = "failed to create tracker"
	msgTrackerAlreadyExists  = "tracker already exists"

	// Success messages
	msgUserCreatedSuccessfully    = "user created successfully"
	msgTrackerCreatedSuccessfully = "tracker created successfully"

	// Allowed sort fields
	allowedUserSortFields    = "id, name, discord_id, mmr, trueskill_mu, trueskill_sigma, created_at, updated_at, trueskill_last_updated"
	allowedTrackerSortFields = "id, discord_id, calculated_mmr, valid, created_at, updated_at, last_updated, ones_current_season_peak, twos_current_season_peak, threes_current_season_peak"
)
