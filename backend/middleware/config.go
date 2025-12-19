package middleware

const (
	// Keycloak server and realm configuration
	keycloakURL = "http://localhost:8080/realms/users" // Keycloak realm URL

	// JWT-specific endpoints
	jwksURL = keycloakURL + "/protocol/openid-connect/certs" // JWKS endpoint for JWT verification

	// Token Introspection-specific configuration
	keycloakIntrospectURL = keycloakURL + "/protocol/openid-connect/token/introspect" // Token introspection endpoint
	serverClientID        = "36601c4e-2027-41f9-b02e-c6a06e20d171"                    // Replace with actual client ID from Keycloak
	serverClientSecret    = "TPuXmlD9X4nzLb5toUTi6MnmWUtoT88U"                        // Replace with actual client secret from Keycloak
)
