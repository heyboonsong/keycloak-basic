# Keycloak Authentication Example - Go Backend

This is a learning example that demonstrates two different methods of validating Keycloak access tokens in a Go backend using **pure JWT verification** (without OIDC library).

## 🎯 Authentication Methods

### 1. JWT Validation (Offline)
**Endpoint:** `/api/todos/private/jwt`

**How it works:**
- Fetches public keys (JWKS) from Keycloak at startup
- Validates the token's RSA signature using the public key
- Verifies token claims (issuer, expiration, audience) locally
- No network call to Keycloak during validation
- **Uses golang-jwt library for pure JWT verification**

**Pros:**
- ⚡ Fast performance (no network latency)
- 🔌 Works offline (no Keycloak dependency after startup)
- 📈 Highly scalable
- 👁️ Transparent - you can see exactly how JWT verification works

**Cons:**
- ⚠️ Cannot detect revoked tokens immediately
- ⏱️ Relies on token expiration only

**Best for:** High-performance applications where tokens are short-lived

---

### 2. Token Introspection (Online)
**Endpoint:** `/api/todos/private/token-introspect`

**How it works:**
- Sends token to Keycloak's introspection endpoint on each request
- Keycloak checks if token is active and not revoked
- Returns detailed token information including username

**Pros:**
- 🔒 Most secure (real-time validation)
- ✅ Detects revoked tokens immediately
- 📊 Gets current token status from Keycloak

**Cons:**
- 🐌 Slower (network call on each request)
- 🌐 Requires Keycloak to be always available
- 📉 Higher load on Keycloak server

**Best for:** Security-critical applications where token revocation matters

---

## 📡 API Endpoints

### Public Endpoint
```
GET /api/todos/public
```
No authentication required

### Protected Endpoints (Authentication Required)
```
GET /api/todos/private/jwt                  (JWT validation)
GET /api/todos/private/token-introspect     (Token introspection)
GET /api/todos/private                      (Default: JWT validation)
```

All configuration is centralized in `middleware/config.go`:

```go
const (
    // Keycloak server and realm configuration
    keycloakURL = "http://localhost:8080/realms/users"
    
    // JWT-specific endpoints
    jwksURL = keycloakURL + "/protocol/openid-connect/certs"
    
    // Token Introspection-specific configuration
    keycloakIntrospectURL = keycloakURL + "/protocol/openid-connect/token/introspect"
    serverClientID        = "your-client-id"
    serverClientSecret    = "your-client-secret"
)
```

**To get the client secret from Keycloak:**
1. Go to Keycloak Admin Console
2. Navigate to your realm → Clients
3. Select your client
4. Go to "Credentials" tab
5. Copy the "Client Secret"
6. Update `serverClientSecret` in `middleware/config.go  keycloakClientID = "your-client-id"

```

## 🚀 Running the Application

```bash
# Start the server
go run ./main.go
```

Server will start on `http://localhost:9000`

## 📚 Code Structure

```
backend/
├── main.go                      # Main application and routes
├── middleware/
│   ├── config.go               # Shared configuration constants
│   ├── helper.go               # Helper functions (token extraction)
│   ├── jwt.go                  # JWT validation with JWKS
│   └── tokenintrospect.go      # Token introspection middleware
├── go.mod                       # Go module dependencies
└── README.MD                    # This file
```

### Key Files Explained:

**main.go**
- Application entry point
- Route definitions
- Server initialization

**middleware/config.go**
- Central configuration for Keycloak
- Shared constants for all authentication methods

**middleware/helper.go**
- Reusable helper functions
- Token extraction from Authorization header

**middleware/jwt.go**
- PThis example uses **pure JWT verification** (golang-jwt) instead of OIDC library for better learning visibility
4. You can see exactly how:
   - JWKS (JSON Web Key Set) is fetched from Keycloak
   - JWK is converted to RSA public key
   - JWT signature is verified using RSA
   - Token claims are validated
5. Choose authentication method based on your security requirements vs performance needs
6. In production, consider caching introspection results with a short TTL for balance

## 🔍 How JWT Verification Works

1. **Startup**: Fetch public keys (JWKS) from Keycloak's `/certs` endpoint
2. **Parse Token**: Extract Key ID (kid) from JWT header
3. **Verify Signature**: Use corresponding RSA public key to verify token signature
4. **Validate Claims**: Check issuer, expiration, and audience
5. **Extract User Info**: Get user details from token claims

## 🔗 Related Documentation

- [Keycloak Token Introspection](https://www.keycloak.org/docs/latest/securing_apps/#_token_introspection_endpoint)
- [JWT Specification](https://datatracker.ietf.org/doc/html/rfc7519)
- [JWKS Specification](https://datatracker.ietf.org/doc/html/rfc7517)
- [golang-jwt Library](https://github.com/golang-jwt/jwt

## 🔍 Key Differences Summary

| Feature | JWT Validation | Token Introspection |
|---------|---------------|-------------------|
| Speed | Fast ⚡ | Slower 🐌 |
| Security | Good 👍 | Excellent 🔒 |
| Revocation Detection | Delayed ⏱️ | Immediate ✅ |
| Network Dependency | Startup only | Every request |
| Keycloak Load | Low | Higher |
| Use Case | High-traffic apps | Security-critical apps |

## 💡 Learning Points

1. **JWT Validation** is like checking an ID card - you verify it's genuine without calling anyone
2. **Token Introspection** is like calling the ID issuer to confirm the ID is still valid
3. Choose based on your security requirements vs performance needs
4. In production, consider caching introspection results with a short TTL for balance

## 🔗 Related Documentation

- [Keycloak Token Introspection](https://www.keycloak.org/docs/latest/securing_apps/#_token_introspection_endpoint)
- [JWT Validation](https://openid.net/specs/openid-connect-core-1_0.html#TokenResponse)
- [OIDC Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html)