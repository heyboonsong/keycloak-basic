package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ===========================
// Token Introspection Middleware
// ===========================
//
// Configuration constants are defined in config.go

// IntrospectionResponse represents Keycloak's token introspection response
type IntrospectionResponse struct {
	Active    bool   `json:"active"`     // Whether the token is active
	Exp       int64  `json:"exp"`        // Token expiration timestamp
	Iat       int64  `json:"iat"`        // Token issued at timestamp
	ClientID  string `json:"client_id"`  // Client ID
	Username  string `json:"username"`   // Username
	TokenType string `json:"token_type"` // Token type (Bearer)
	Sub       string `json:"sub"`        // Subject (user ID)
	Email     string `json:"email"`      // User email
}

// TokenIntrospectionMiddleware validates access tokens using Keycloak's introspection endpoint
// This method:
// - Sends token to Keycloak for validation on each request
// - Checks if token is revoked or blacklisted
// - More secure but slower (network call to Keycloak)
// - Good for security-critical scenarios
func TokenIntrospectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Bearer token from Authorization header
		token, ok := extractBearerToken(c)
		if !ok {
			return // Error response already sent by extractBearerToken
		}

		// Introspect the token with Keycloak
		introspectResp, err := introspectToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Failed to introspect token",
				"detail": err.Error(),
				"method": "Token Introspection",
			})
			c.Abort()
			return
		}

		// Step 4: Check if token is active
		if !introspectResp.Active {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token is not active or has been revoked",
				"method": "Token Introspection",
			})
			c.Abort()
			return
		}

		// Store user info in context for use in handlers
		c.Set("user_id", introspectResp.Sub)
		c.Set("email", introspectResp.Email)

		// Token is valid and active, continue to the next handler
		c.Next()
	}
}

// introspectToken calls Keycloak's introspection endpoint to validate the token
func introspectToken(token string) (*IntrospectionResponse, error) {
	// Prepare form data for introspection request
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", serverClientID)
	data.Set("client_secret", serverClientSecret)

	// Create HTTP request to Keycloak introspection endpoint
	req, err := http.NewRequest("POST", keycloakIntrospectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create introspection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request to Keycloak
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read introspection response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var introspectResp IntrospectionResponse
	if err := json.Unmarshal(body, &introspectResp); err != nil {
		return nil, fmt.Errorf("failed to parse introspection response: %w", err)
	}

	return &introspectResp, nil
}
