package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/supabase-community/supabase-go"
)

// DiscordAuth provides unified Discord OAuth authentication via Supabase Auth
// This replaces both the main app auth and USL auth systems
type DiscordAuth struct {
	supabaseClient  *supabase.Client
	adminDiscordIDs []string
	supabaseURL     string // Internal Supabase URL for backend calls
	publicURL       string // Public Supabase URL for OAuth redirects
	anonKey         string // Anon key for client-side operations
}

func NewDiscordAuth(supabaseClient *supabase.Client, adminDiscordIDs []string, supabaseURL, publicURL, anonKey string) *DiscordAuth {
	return &DiscordAuth{
		supabaseClient:  supabaseClient,
		adminDiscordIDs: adminDiscordIDs,
		supabaseURL:     supabaseURL,
		publicURL:       publicURL,
		anonKey:         anonKey,
	}
}

func (a *DiscordAuth) LoginForm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if a.IsAuthenticated(r) {
		if strings.HasPrefix(r.URL.Path, "/usl") {
			http.Redirect(w, r, "/usl/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/users", http.StatusSeeOther)
		}
		return
	}

	var redirectTo string
	if strings.HasPrefix(r.URL.Path, "/usl") {
		redirectTo = "http://127.0.0.1:8080/auth/callback?redirect=usl"
	} else {
		redirectTo = "http://127.0.0.1:8080/auth/callback?redirect=main"
	}

	discordOAuthURL := fmt.Sprintf("%s/auth/v1/authorize?provider=discord&redirect_to=%s",
		a.publicURL, redirectTo)

	isUSL := strings.HasPrefix(r.URL.Path, "/usl")
	var title, heading, infoText string

	if isUSL {
		title = "USL Admin Login"
		heading = "USL Administration"
		infoText = "Sign in with Discord to access the USL management system."
	} else {
		title = "Sign In"
		heading = "Sign In"
		infoText = "Sign in with Discord to access the application."
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .login-container {
            background-color: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            max-width: 400px;
            width: 100%%;
            text-align: center;
        }
        h2 {
            color: #333;
            margin-bottom: 20px;
        }
        .usl-header {
            color: #1a365d;
            margin-bottom: 10px;
            font-size: 24px;
            font-weight: bold;
        }
        .discord-btn {
            display: block;
            width: 100%%;
            padding: 15px;
            background-color: #5865F2;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-size: 16px;
            font-weight: bold;
            transition: background-color 0.3s;
            text-align: center;
            box-sizing: border-box;
        }
        .discord-btn:hover {
            background-color: #4752C4;
        }
        .info {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .error {
            color: red;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h2 class="%s">%s</h2>
        <div class="info">%s</div>
        %s
        <a href="%s" class="discord-btn">
            🎮 Sign in with Discord
        </a>
    </div>
</body>
</html>
`, title, func() string {
		if isUSL {
			return "usl-header"
		}
		return ""
	}(), heading, infoText, a.getErrorMessage(r), discordOAuthURL)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write auth error response: %v", err)
	}
}

func (a *DiscordAuth) AuthCallback(w http.ResponseWriter, r *http.Request) {
	redirectParam := r.URL.Query().Get("redirect")
	var finalRedirect string

	switch redirectParam {
	case "usl":
		finalRedirect = "/usl/admin"
	case "main":
		finalRedirect = "/users"
	default:
		finalRedirect = "/users" // Default to main app
	}

	// HTML page that extracts tokens from URL fragment and sets session
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Processing Login...</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background-color: #f5f5f5;
        }
        .processing {
            text-align: center;
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
    </style>
</head>
<body>
    <div class="processing">
        <h2>Processing Login...</h2>
        <p>Please wait while we complete your authentication.</p>
    </div>
    <script>
        // Extract tokens from URL fragment (Supabase returns them in hash)
        const hash = window.location.hash.substring(1);
        const params = new URLSearchParams(hash);
        const accessToken = params.get('access_token');
        const refreshToken = params.get('refresh_token');
        
        if (accessToken) {
            // Send tokens to server for validation and session setup
            fetch('/auth/process', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    access_token: accessToken,
                    refresh_token: refreshToken 
                })
            }).then(response => {
                if (response.ok) {
                    // Successful authentication - redirect to final destination
                    window.location.href = '%s';
                } else {
                    // Authentication failed
                    window.location.href = '/login?error=unauthorized';
                }
            }).catch(error => {
                console.error('Auth error:', error);
                window.location.href = '/login?error=invalid';
            });
        } else {
            // No token found - redirect to login with error
            window.location.href = '/login?error=invalid';
        }
    </script>
</body>
</html>
`, finalRedirect)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(html)); err != nil {
		log.Printf("Failed to write auth callback response: %v", err)
	}
}

// ProcessTokens handles the access token validation and session setup
func (a *DiscordAuth) ProcessTokens(w http.ResponseWriter, r *http.Request) {
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

	user, err := a.validateTokensAndGetUser(req.AccessToken)
	if err != nil {
		log.Printf("Error validating tokens: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	if !a.isUserAuthorized(user) {
		log.Printf("User not authorized: %+v", user)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	a.setSessionCookies(w, req.AccessToken, req.RefreshToken)

	w.WriteHeader(http.StatusOK)
}

func (a *DiscordAuth) IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("auth_access_token")
	if err != nil {
		return false
	}

	// Validate token and get user info
	user, err := a.validateTokensAndGetUser(cookie.Value)
	if err != nil {
		return false
	}

	return a.isUserAuthorized(user)
}

// RequireAuth middleware that protects routes with Discord OAuth
// This replaces both the main app auth middleware and USL auth middleware
func (a *DiscordAuth) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.IsAuthenticated(r) {
			if strings.HasPrefix(r.URL.Path, "/usl") {
				http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}
		next(w, r)
	}
}

func (a *DiscordAuth) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to appropriate login
	if strings.HasPrefix(r.URL.Path, "/usl") {
		http.Redirect(w, r, "/usl/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Helper methods

func (a *DiscordAuth) validateTokensAndGetUser(accessToken string) (map[string]interface{}, error) {
	// Use the main Supabase client with the user's access token

	userClient := a.supabaseClient.Auth.WithToken(accessToken)

	user, err := userClient.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Supabase: %w", err)
	}

	// Extract user information
	userMap := map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"user_metadata": user.UserMetadata,
		"app_metadata":  user.AppMetadata,
	}

	return userMap, nil
}

func (a *DiscordAuth) isUserAuthorized(user map[string]interface{}) bool {
	// Extract Discord ID from user metadata
	// The exact location depends on how Supabase stores Discord user info
	discordID := ""

	// Debug: Log the entire user object to see what we're working with

	// Try to get Discord ID from various possible locations
	if metadata, ok := user["user_metadata"].(map[string]interface{}); ok {
		if id, exists := metadata["provider_id"].(string); exists {
			discordID = id
		} else if id, exists := metadata["sub"].(string); exists {
			discordID = id
		} else if id, exists := metadata["discord_id"].(string); exists {
			discordID = id
		}
	}

	if discordID == "" {
		return false
	}

	for _, adminID := range a.adminDiscordIDs {
		if discordID == adminID {
			return true
		}
	}

	return false
}

func (a *DiscordAuth) setSessionCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // 1 hour
	})

	if refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400 * 7, // 7 days
		})
	}
}

func (a *DiscordAuth) getErrorMessage(r *http.Request) string {
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
