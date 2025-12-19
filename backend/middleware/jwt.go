package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ===========================
// JWT Verification for Keycloak
// ===========================
//
// This file provides JWT token verification for Keycloak authentication
// WITHOUT using the OIDC library. This demonstrates how JWT validation works:
//
// 1. Fetch Public Keys (JWKS) from Keycloak at startup
//    - Keycloak exposes public keys at /protocol/openid-connect/certs
//    - These keys are used to verify token signatures
//
// 2. Parse and Verify JWT Token
//    - Extract the Key ID (kid) from the token header
//    - Find the matching public key
//    - Verify the RSA signature using the public key
//    - Validate claims (issuer, expiration, audience)
//
// 3. Return Token Claims
//    - If valid, return user information from the token
//    - If invalid, return error
//
// This approach is:
// ✅ Simple and transparent - you can see exactly how JWT works
// ✅ Fast - no network calls during validation (keys cached)
// ✅ Standard - uses only JWT and RSA cryptography
//
// Configuration constants are defined in config.go

// ===========================
// Data Structures
// ===========================

// JWKS represents the JSON Web Key Set from Keycloak
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key
type JWK struct {
	Kid string `json:"kid"` // Key ID
	Kty string `json:"kty"` // Key Type (RSA)
	Alg string `json:"alg"` // Algorithm (RS256)
	Use string `json:"use"` // Usage (sig for signature)
	N   string `json:"n"`   // Modulus (public key)
	E   string `json:"e"`   // Exponent (public key)
}

// CustomClaims represents JWT claims structure
type CustomClaims struct {
	jwt.RegisteredClaims
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
}

// ===========================
// Global Variables
// ===========================

var publicKeys map[string]*rsa.PublicKey

// ===========================
// Initialization
// ===========================

// GetKeycloakPublicKey fetches public keys from Keycloak for JWT verification
func GetKeycloakPublicKey() error {
	// Fetch JWKS (JSON Web Key Set) from Keycloak
	keys, err := fetchJWKS()
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Convert JWKs to RSA public keys
	publicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range keys.Keys {
		pubKey, err := jwkToRSAPublicKey(key)
		if err != nil {
			return fmt.Errorf("failed to convert JWK to RSA public key: %w", err)
		}
		publicKeys[key.Kid] = pubKey
	}

	return nil
}

// ===========================
// JWT Verification
// ===========================

// verifyToken validates a JWT token using Keycloak's public keys
// Returns the parsed claims if valid, error otherwise
func verifyToken(tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found in token")
		}

		// Get the corresponding public key
		publicKey, exists := publicKeys[kid]
		if !exists {
			return nil, fmt.Errorf("public key not found for kid: %s", kid)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate issuer
	expectedIssuer := keycloakURL
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", expectedIssuer, claims.Issuer)
	}

	// Validate expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Validate audience (optional, based on your needs)
	if len(claims.Audience) > 0 {
		validAudience := false
		for _, aud := range claims.Audience {
			if aud == "account" {
				validAudience = true
				break
			}
		}
		if !validAudience {
			return nil, fmt.Errorf("invalid audience")
		}
	}

	return claims, nil
}

// ===========================
// Helper Functions
// ===========================

// fetchJWKS fetches the JSON Web Key Set from Keycloak
func fetchJWKS() (*JWKS, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	return &jwks, nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert bytes to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}

// ===========================
// JWT Validation Middleware
// ===========================

// JWTAuthMiddleware validates access tokens using JWT verification
// This method:
// - Validates the token signature using RSA public keys from Keycloak
// - Checks token expiration, issuer, and audience claims
// - Faster - no network call during validation (keys fetched at startup)
// - Good for high-performance scenarios
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Bearer token from Authorization header
		token, ok := extractBearerToken(c)
		if !ok {
			return // Error response already sent by extractBearerToken
		}

		// Verify the JWT token
		// This validates:
		// - Token signature using RSA public key
		// - Token expiration time
		// - Token issuer matches Keycloak
		// - Token audience (if present)
		claims, err := verifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid or expired JWT token",
				"detail": err.Error(),
				"method": "JWT Validation",
			})
			c.Abort()
			return
		}

		// Store user info in context for use in handlers
		c.Set("user_id", claims.Subject)
		c.Set("email", claims.Email)

		// Token is valid, continue to the next handler
		c.Next()
	}
}
