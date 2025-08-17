package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/supabase-community/supabase-go"
	"usl-server/internal/config"
	"usl-server/internal/templates"
)

const (
	metadataProviderID = "provider_id"
	metadataSub        = "sub"
	metadataDiscordID  = "discord_id"
	userMetadataKey    = "user_metadata"

	// Route constants
	uslAdminRoute = "/usl/admin"
	usersRoute    = "/users"
	uslLoginRoute = "/usl/login"
	loginRoute    = "/login"

	// Cookie constants
	accessTokenCookie  = "auth_access_token"
	refreshTokenCookie = "auth_refresh_token"

	// URL prefixes
	uslPrefix = "/usl"
)

type DiscordAuth struct {
	supabaseClient  *supabase.Client
	adminDiscordIDs []string
	supabaseURL     string
	publicURL       string
	anonKey         string
	envConfig       *config.EnvironmentConfig
}

func NewDiscordAuth(supabaseClient *supabase.Client, adminDiscordIDs []string, supabaseURL, publicURL, anonKey string, envConfig config.EnvironmentConfig) *DiscordAuth {
	return &DiscordAuth{
		supabaseClient:  supabaseClient,
		adminDiscordIDs: adminDiscordIDs,
		supabaseURL:     supabaseURL,
		publicURL:       publicURL,
		anonKey:         anonKey,
		envConfig:       &envConfig,
	}
}

func (auth *DiscordAuth) LoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if auth.IsAuthenticated(r) {
		auth.redirectAuthenticated(w, r)
		return
	}

	appBaseURL := auth.getAppBaseURL()
	redirectTo := auth.buildRedirectURL(appBaseURL, r.URL.Path)
	discordOAuthURL := fmt.Sprintf("%s/auth/v1/authorize?provider=discord&redirect_to=%s",
		auth.publicURL, redirectTo)

	title, heading, infoText := auth.getLoginPageContent(r.URL.Path)

	html := fmt.Sprintf(templates.LoginFormHTML,
		title,
		auth.getHeaderClass(r.URL.Path),
		heading,
		infoText,
		auth.getErrorMessage(r),
		discordOAuthURL)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write auth error response: %v", err)
	}
}

func (auth *DiscordAuth) AuthCallback(w http.ResponseWriter, r *http.Request) {
	redirectParam := r.URL.Query().Get("redirect")
	var finalRedirect string

	switch redirectParam {
	case "usl":
		finalRedirect = uslAdminRoute
	case "main":
		finalRedirect = usersRoute
	default:
		finalRedirect = usersRoute
	}

	html := fmt.Sprintf(templates.AuthCallbackHTML, finalRedirect)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write auth callback response: %v", err)
	}
}

// ProcessTokens handles the access token validation and session setup
func (auth *DiscordAuth) ProcessTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding token request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.validateTokensAndGetUser(req.AccessToken)
	if err != nil {
		log.Printf("Error validating tokens: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	if !auth.isUserAuthorized(user) {
		log.Printf("User not authorized: %+v", user)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	auth.setSessionCookies(w, req.AccessToken, req.RefreshToken)

	w.WriteHeader(http.StatusOK)
}

func (auth *DiscordAuth) IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(accessTokenCookie)
	if err != nil {
		return false
	}

	// Validate token and get user info
	user, err := auth.validateTokensAndGetUser(cookie.Value)
	if err != nil {
		return false
	}

	return auth.isUserAuthorized(user)
}

func (auth *DiscordAuth) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAuthenticated(r) {
			if auth.isUSLPath(r.URL.Path) {
				http.Redirect(w, r, uslLoginRoute, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, loginRoute, http.StatusSeeOther)
			}
			return
		}
		next(w, r)
	}
}

func (auth *DiscordAuth) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	if auth.isUSLPath(r.URL.Path) {
		http.Redirect(w, r, uslLoginRoute, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
	}
}

func (auth *DiscordAuth) validateTokensAndGetUser(accessToken string) (map[string]interface{}, error) {
	userClient := auth.supabaseClient.Auth.WithToken(accessToken)

	user, err := userClient.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Supabase: %w", err)
	}

	userInfo := map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"user_metadata": user.UserMetadata,
		"app_metadata":  user.AppMetadata,
	}

	return userInfo, nil
}

func (auth *DiscordAuth) isUserAuthorized(user map[string]interface{}) bool {
	discordID := auth.extractDiscordID(user)
	if discordID == "" {
		return false
	}

	for _, adminID := range auth.adminDiscordIDs {
		if discordID == adminID {
			return true
		}
	}

	return false
}

func (auth *DiscordAuth) extractDiscordID(user map[string]interface{}) string {
	metadata, ok := user[userMetadataKey].(map[string]interface{})
	if !ok {
		return ""
	}

	fields := []string{metadataProviderID, metadataSub, metadataDiscordID}
	for _, field := range fields {
		if id, exists := metadata[field].(string); exists && id != "" {
			return id
		}
	}

	return ""
}

func (auth *DiscordAuth) setSessionCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	if refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     refreshTokenCookie,
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400 * 7,
		})
	}
}

func (auth *DiscordAuth) getErrorMessage(r *http.Request) string {
	errorType := r.URL.Query().Get("error")
	switch errorType {
	case "invalid":
		return `<div class="error">Login failed. Please try again.</div>`
	case "unauthorized":
		return `<div class="error">Access denied. Your Discord account is not authorized.</div>`
	default:
		return ""
	}
}

func (auth *DiscordAuth) getAppBaseURL() string {
	if auth.envConfig != nil && auth.envConfig.AppBaseURL != "" {
		return auth.envConfig.AppBaseURL
	}

	return "http://localhost:8080"
}

func (auth *DiscordAuth) GetAppBaseURL() string {
	return auth.getAppBaseURL()
}

func (auth *DiscordAuth) isUSLPath(path string) bool {
	return strings.HasPrefix(path, uslPrefix)
}

func (auth *DiscordAuth) redirectAuthenticated(w http.ResponseWriter, r *http.Request) {
	if auth.isUSLPath(r.URL.Path) {
		http.Redirect(w, r, uslAdminRoute, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, usersRoute, http.StatusSeeOther)
	}
}

func (auth *DiscordAuth) buildRedirectURL(appBaseURL, path string) string {
	if auth.isUSLPath(path) {
		return fmt.Sprintf("%s/auth/callback?redirect=usl", appBaseURL)
	}
	return fmt.Sprintf("%s/auth/callback?redirect=main", appBaseURL)
}

func (auth *DiscordAuth) getLoginPageContent(path string) (title, heading, infoText string) {
	if auth.isUSLPath(path) {
		return "USL Admin Login", "USL Administration", "Sign in with Discord to access the USL management system."
	}
	return "Sign In", "Sign In", "Sign in with Discord to access the application."
}

func (auth *DiscordAuth) getHeaderClass(path string) string {
	if auth.isUSLPath(path) {
		return "usl-header"
	}
	return ""
}
