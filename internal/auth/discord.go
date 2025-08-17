package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"usl-server/internal/config"
	"usl-server/internal/templates"

	"github.com/supabase-community/supabase-go"
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
	// RENDER DEBUG: Log all incoming request details at OAuth callback
	log.Printf("RENDER DEBUG: OAuth callback received")
	log.Printf("RENDER DEBUG: Request URL: %s", r.URL.String())
	log.Printf("RENDER DEBUG: Request Method: %s", r.Method)
	log.Printf("RENDER DEBUG: Request Host: %s", r.Host)
	log.Printf("RENDER DEBUG: Request RemoteAddr: %s", r.RemoteAddr)
	log.Printf("RENDER DEBUG: Request Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("RENDER DEBUG:   %s: %s", name, value)
		}
	}

	// Log all query parameters
	log.Printf("RENDER DEBUG: Query Parameters:")
	for name, values := range r.URL.Query() {
		for _, value := range values {
			if name == "access_token" && len(value) > 50 {
				log.Printf("RENDER DEBUG:   %s: %s... (truncated)", name, value[:50])
			} else {
				log.Printf("RENDER DEBUG:   %s: %s", name, value)
			}
		}
	}

	// Log fragment from URL if present
	if r.URL.Fragment != "" {
		log.Printf("RENDER DEBUG: URL Fragment: %s", r.URL.Fragment)
	}

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

	log.Printf("RENDER DEBUG: Final redirect destination: %s", finalRedirect)

	html := fmt.Sprintf(templates.AuthCallbackHTML, finalRedirect)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write auth callback response: %v", err)
	}
}

// ProcessTokens handles the access token validation and session setup
func (auth *DiscordAuth) ProcessTokens(w http.ResponseWriter, r *http.Request) {
	// RENDER DEBUG: Log incoming token processing request
	log.Printf("RENDER DEBUG: ProcessTokens called")
	log.Printf("RENDER DEBUG: Request Method: %s", r.Method)
	log.Printf("RENDER DEBUG: Request Host: %s", r.Host)
	log.Printf("RENDER DEBUG: Request RemoteAddr: %s", r.RemoteAddr)
	log.Printf("RENDER DEBUG: Request Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("RENDER DEBUG:   %s: %s", name, value)
		}
	}

	if r.Method != http.MethodPost {
		log.Printf("RENDER DEBUG: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("RENDER DEBUG: Error decoding token request: %v", err)
		log.Printf("Error decoding token request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Debug logging for token validation
	log.Printf("RENDER DEBUG: Received access token length: %d", len(req.AccessToken))
	log.Printf("RENDER DEBUG: Access token (first 20 chars): %s...", req.AccessToken[:min(20, len(req.AccessToken))])
	if req.RefreshToken != "" {
		log.Printf("RENDER DEBUG: Received refresh token length: %d", len(req.RefreshToken))
		log.Printf("RENDER DEBUG: Refresh token (first 20 chars): %s...", req.RefreshToken[:min(20, len(req.RefreshToken))])
	}
	log.Printf("DEBUG: Token validation starting")

	user, err := auth.validateTokensAndGetUser(req.AccessToken)
	if err != nil {
		log.Printf("RENDER DEBUG: Token validation FAILED in production environment")
		log.Printf("Error validating tokens: %v", err)
		log.Printf("DEBUG: Token validation failed - full error: %+v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	log.Printf("RENDER DEBUG: Token validation SUCCEEDED in production environment")
	log.Printf("DEBUG: Token validation successful, user ID: %v", user["id"])

	if !auth.isUserAuthorized(user) {
		log.Printf("RENDER DEBUG: User authorization FAILED")
		log.Printf("User not authorized: %+v", user)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	log.Printf("RENDER DEBUG: User authorization SUCCEEDED, setting session cookies")
	auth.setSessionCookies(w, req.AccessToken, req.RefreshToken)

	log.Printf("RENDER DEBUG: ProcessTokens completed successfully")
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
	log.Printf("DEBUG: Creating anon client with URL: %s", auth.supabaseURL)
	log.Printf("DEBUG: Using anon key (first 20 chars): %s...", auth.anonKey[:min(20, len(auth.anonKey))])

	anonClient, err := supabase.NewClient(auth.supabaseURL, auth.anonKey, nil)
	if err != nil {
		log.Printf("DEBUG: Failed to create anon client: %v", err)
		return nil, fmt.Errorf("failed to create anon client: %w", err)
	}
	log.Printf("DEBUG: Anon client created successfully")

	log.Printf("DEBUG: Creating user client with WithToken")
	userClient := anonClient.Auth.WithToken(accessToken)
	log.Printf("DEBUG: User client created, calling GetUser()")

	user, err := userClient.GetUser()
	if err != nil {
		log.Printf("DEBUG: GetUser() failed with error: %v", err)
		log.Printf("DEBUG: Error type: %T", err)
		return nil, fmt.Errorf("failed to get user from Supabase: %w", err)
	}
	log.Printf("DEBUG: GetUser() succeeded, user email: %s", user.Email)

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
