package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ===========================
// Helper Functions
// ===========================

// extractBearerToken extracts the Bearer token from the Authorization header
// Returns the token string and a boolean indicating success
func extractBearerToken(c *gin.Context) (string, bool) {
	// Step 1: Check if Authorization header exists
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header required",
		})
		c.Abort()
		return "", false
	}

	// Step 2: Extract Bearer token from "Bearer <token>" format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format. Expected: Bearer <token>",
		})
		c.Abort()
		return "", false
	}

	return parts[1], true
}
