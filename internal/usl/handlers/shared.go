package handlers

import (
	"usl-server/internal/usl"
)

// Type aliases for cleaner code
type USLUser = usl.USLUser
type USLUserTracker = usl.USLUserTracker

const (
	// USL Configuration
	USLDiscordGuildID = "1390537743385231451"

	// Rocket League Business Rules
	MinMMR             = 0
	MaxMMR             = 3000  // SSL is around 1900-2000, allow buffer for edge cases
	MaxGames           = 10000 // Reasonable season game limit
	MinDiscordIDLength = 17
	MaxDiscordIDLength = 19

	// Validation Error Codes
	ValidationCodeRequired      = "required"
	ValidationCodeInvalidFormat = "invalid_format"
	ValidationCodeOutOfRange    = "out_of_range"
	ValidationCodeLogicalError  = "logical_error"
	ValidationCodeInvalidURL    = "invalid_url"
	ValidationCodeNoData        = "no_data"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

type FormField string

const (
	// Basic Entity Fields
	FormFieldID        FormField = "id"
	FormFieldDiscordID FormField = "discord_id"
	FormFieldURL       FormField = "url"
	FormFieldName      FormField = "name"
	FormFieldActive    FormField = "active"
	FormFieldBanned    FormField = "banned"
	FormFieldValid     FormField = "valid"

	// 1v1 MMR Fields
	FormFieldOnesCurrentPeak   FormField = "ones_current_peak"
	FormFieldOnesPreviousPeak  FormField = "ones_previous_peak"
	FormFieldOnesAllTimePeak   FormField = "ones_all_time_peak"
	FormFieldOnesCurrentGames  FormField = "ones_current_games"
	FormFieldOnesPreviousGames FormField = "ones_previous_games"

	// 2v2 MMR Fields
	FormFieldTwosCurrentPeak   FormField = "twos_current_peak"
	FormFieldTwosPreviousPeak  FormField = "twos_previous_peak"
	FormFieldTwosAllTimePeak   FormField = "twos_all_time_peak"
	FormFieldTwosCurrentGames  FormField = "twos_current_games"
	FormFieldTwosPreviousGames FormField = "twos_previous_games"

	// 3v3 MMR Fields
	FormFieldThreesCurrentPeak   FormField = "threes_current_peak"
	FormFieldThreesPreviousPeak  FormField = "threes_previous_peak"
	FormFieldThreesAllTimePeak   FormField = "threes_all_time_peak"
	FormFieldThreesCurrentGames  FormField = "threes_current_games"
	FormFieldThreesPreviousGames FormField = "threes_previous_games"
)

// TemplateName represents typed template names
type TemplateName string

const (
	TemplateUSLUsers          TemplateName = "users-list-page"
	TemplateUSLUsersTable     TemplateName = "users-table-fragment"
	TemplateUSLUserDetail     TemplateName = "user-detail-page"
	TemplateUSLTrackers       TemplateName = "trackers-list-page"
	TemplateUSLTrackerDetail  TemplateName = "tracker-detail-page"
	TemplateUSLTrackerNew     TemplateName = "tracker-new-page"
	TemplateUSLTrackerEdit    TemplateName = "tracker-edit-page"
	TemplateUSLAdminDashboard TemplateName = "admin-dashboard-page"
)

// Template Data Types for consistent structure and better maintainability

type BasePageData struct {
	Title       string
	CurrentPage string
}

type SearchConfig struct {
	SearchPlaceholder string
	SearchURL         string
	SearchTarget      string
	ClearURL          string
	ShowFilters       bool
	Query             string
	StatusFilter      string
}

type TrackerFormData struct {
	BasePageData
	Tracker *USLUserTracker
	User    *USLUser
	Errors  map[string]string
}

type UserDetailData struct {
	BasePageData
	User         *USLUser
	UserTrackers []*USLUserTracker
}

type TrackerDetailData struct {
	BasePageData
	Tracker *USLUserTracker
	User    *USLUser
}

type UsersListData struct {
	BasePageData
	Users        []*USLUser
	SearchConfig SearchConfig
}

type TrackersListData struct {
	BasePageData
	Trackers     []*USLUserTracker
	SearchConfig SearchConfig
}

type AdminDashboardData struct {
	BasePageData
	Stats struct {
		TotalUsers    int `json:"total_users"`
		ActiveUsers   int `json:"active_users"`
		TotalTrackers int `json:"total_trackers"`
		ValidTrackers int `json:"valid_trackers"`
	}
}

type TrueSkillUpdateResult struct {
	Success         bool   `json:"success"`
	HadTrackers     bool   `json:"hadTrackers"`
	Error           string `json:"error,omitempty"`
	TrueSkillResult *struct {
		Mu    float64 `json:"mu"`
		Sigma float64 `json:"sigma"`
	} `json:"trueSkillResult,omitempty"`
	UserName string `json:"userName"`
}
