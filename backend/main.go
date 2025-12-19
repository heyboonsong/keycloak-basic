package main

import (
	"fmt"
	"keycloak-basic-backend/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Initializing JWT verification...")
	if err := middleware.GetKeycloakPublicKey(); err != nil {
		panic(fmt.Sprintf("Failed to initialize JWT: %v", err))
	}
	fmt.Println("✓ JWT initialized successfully")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/api/todos/public", getTodosPublic)

	// Method 1: JWT Validation (Fast, Offline)
	// Validates token cryptographically without calling Keycloak
	// ✅ Faster performance (no network call)
	// ✅ Works offline
	// ❌ Cannot detect revoked tokens immediately
	router.GET("/api/todos/private/jwt", middleware.JWTAuthMiddleware(), getTodosPrivate)

	// Method 2: Token Introspection (Secure, Online)
	// Validates token by calling Keycloak introspection endpoint
	// ✅ Detects revoked tokens immediately
	// ✅ More secure
	// ❌ Slower (requires network call to Keycloak)
	router.GET("/api/todos/private/token-introspect", middleware.TokenIntrospectionMiddleware(), getTodosPrivate)

	serverPort := ":9000"
	fmt.Printf("🚀 Server starting on %s\n", serverPort)

	if err := router.Run(serverPort); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var todos = []Todo{
	{
		ID:          "1",
		Title:       "Learn Go",
		Description: "Study Go programming language",
		Completed:   false,
	},
	{
		ID:          "2",
		Title:       "Build API",
		Description: "Create REST API with Gin",
		Completed:   false,
	},
	{
		ID:          "3",
		Title:       "Integrate Keycloak",
		Description: "Add authentication with Keycloak",
		Completed:   true,
	},
}

func getTodosPublic(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Public endpoint - no authentication required",
		"data":    todos,
	})
}

func getTodosPrivate(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")

	response := gin.H{
		"authenticated_user": gin.H{
			"user_id": userID,
			"email":   email,
		},
		"message": "Protected endpoint - authentication required",
		"data":    todos,
	}

	c.JSON(http.StatusOK, response)
}
