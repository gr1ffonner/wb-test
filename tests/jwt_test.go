package tests

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"wb-test/pkg/utils/jwt"
)

func TestJWTFlow(t *testing.T) {
	// JWT Token Generation and Validation Example
	demonstrateJWTFlow(t)
}

// demonstrateJWTFlow shows token generation and validation
func demonstrateJWTFlow(t *testing.T) {
	log := slog.Default()
	log.Info("=== JWT Token Flow Demo ===")

	// 1. Generate a JWT token
	userID := 123
	username := "john_doe"

	token, err := jwt.GenerateJWT(userID, username)
	require.NoError(t, err, "Failed to generate JWT")
	log.Info("JWT Token generated", "user_id", userID, "username", username, "token", token[:20]+"...")

	// 2. Validate the token
	err = jwt.ValidateToken(token)
	require.NoError(t, err, "Token validation failed")
	log.Info("Token validation successful")

	// 3. Parse the token
	claims, err := jwt.ParseJWT(token)
	require.NoError(t, err, "Failed to parse JWT")
	assert.Equal(t, userID, claims.UserID, "UserID should match")
	assert.Equal(t, username, claims.Username, "Username should match")
	log.Info("JWT Token validated", "user_id", claims.UserID, "username", claims.Username, "expires_at", claims.ExpiresAt)

	// 4. Demonstrate error handling with invalid token
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"
	_, err = jwt.ParseJWT(invalidToken)
	assert.Error(t, err, "Invalid token should be rejected")
	log.Info("Invalid token correctly rejected", "error", err)

	log.Info("=== JWT Demo Complete ===")
}

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		username string
		wantErr  bool
	}{
		{
			name:     "valid user data",
			userID:   123,
			username: "john_doe",
			wantErr:  false,
		},
		{
			name:     "zero user ID",
			userID:   0,
			username: "anonymous",
			wantErr:  false,
		},
		{
			name:     "empty username",
			userID:   456,
			username: "",
			wantErr:  false,
		},
		{
			name:     "negative user ID",
			userID:   -1,
			username: "negative_user",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.GenerateJWT(tt.userID, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Contains(t, token, ".")
		})
	}
}

func TestParseJWT(t *testing.T) {
	// Generate a valid token first
	userID := 123
	username := "test_user"
	validToken, err := jwt.GenerateJWT(userID, username)
	require.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		wantUserID   int
		wantUsername string
		wantErr      bool
	}{
		{
			name:         "valid token",
			token:        validToken,
			wantUserID:   userID,
			wantUsername: username,
			wantErr:      false,
		},
		{
			name:         "invalid token format",
			token:        "invalid.token.format",
			wantUserID:   0,
			wantUsername: "",
			wantErr:      true,
		},
		{
			name:         "empty token",
			token:        "",
			wantUserID:   0,
			wantUsername: "",
			wantErr:      true,
		},
		{
			name:         "malformed token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			wantUserID:   0,
			wantUsername: "",
			wantErr:      true,
		},
		{
			name:         "token with wrong signature",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoidGVzdF91c2VyIn0.wrong_signature",
			wantUserID:   0,
			wantUsername: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwt.ParseJWT(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, claims)
			assert.Equal(t, tt.wantUserID, claims.UserID)
			assert.Equal(t, tt.wantUsername, claims.Username)
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Generate a valid token
	validToken, err := jwt.GenerateJWT(123, "test_user")
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jwt.ValidateToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    bool
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature",
			wantToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature",
			wantErr:    false,
		},
		{
			name:       "missing authorization header",
			authHeader: "",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "invalid bearer format",
			authHeader: "Basic dGVzdDp0ZXN0",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "bearer without token",
			authHeader: "Bearer ",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "bearer with empty token",
			authHeader: "Bearer",
			wantToken:  "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := jwt.ExtractTokenFromHeader(req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	// Create test claims
	testClaims := &jwt.Claims{
		UserID:   123,
		Username: "test_user",
	}

	tests := []struct {
		name       string
		ctx        context.Context
		wantClaims *jwt.Claims
		wantOK     bool
	}{
		{
			name:       "user in context",
			ctx:        context.WithValue(context.Background(), "user", testClaims),
			wantClaims: testClaims,
			wantOK:     true,
		},
		{
			name:       "no user in context",
			ctx:        context.Background(),
			wantClaims: nil,
			wantOK:     false,
		},
		{
			name:       "wrong type in context",
			ctx:        context.WithValue(context.Background(), "user", "not a claims object"),
			wantClaims: nil,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, ok := jwt.GetUserFromContext(tt.ctx)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantClaims, claims)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// Generate original token
	userID := 123
	username := "test_user"
	originalToken, err := jwt.GenerateJWT(userID, username)
	require.NoError(t, err)

	// Parse original token to get original expiration
	originalClaims, err := jwt.ParseJWT(originalToken)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token refresh",
			token:   originalToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add a longer delay to ensure different timestamps
			time.Sleep(100 * time.Millisecond)

			refreshedToken, err := jwt.RefreshToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, refreshedToken)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, refreshedToken)

			// Parse refreshed token
			refreshedClaims, err := jwt.ParseJWT(refreshedToken)
			require.NoError(t, err)

			// Check that user data is preserved
			assert.Equal(t, originalClaims.UserID, refreshedClaims.UserID)
			assert.Equal(t, originalClaims.Username, refreshedClaims.Username)

			// Check that expiration time is extended (should be later than original)
			expirationTimeDiff := refreshedClaims.ExpiresAt.Time.Sub(originalClaims.ExpiresAt.Time)
			t.Logf("Expiration time difference: %v", expirationTimeDiff)
			assert.True(t, refreshedClaims.ExpiresAt.Time.After(originalClaims.ExpiresAt.Time),
				"Refreshed expiration time should be after original expiration time")

			// Check that issued time is at least the same or later (JWT library might use same timestamp)
			issuedTimeDiff := refreshedClaims.IssuedAt.Time.Sub(originalClaims.IssuedAt.Time)
			t.Logf("Issued time difference: %v", issuedTimeDiff)
			assert.True(t, !refreshedClaims.IssuedAt.Time.Before(originalClaims.IssuedAt.Time),
				"Refreshed issued time should not be before original issued time")

			// Verify that both tokens are valid
			assert.NoError(t, jwt.ValidateToken(originalToken))
			assert.NoError(t, jwt.ValidateToken(refreshedToken))

			// Verify the tokens are different (due to different timestamps)
			assert.NotEqual(t, originalToken, refreshedToken)
		})
	}
}

func TestJWTClaimsStructure(t *testing.T) {
	// Test that claims structure works correctly
	userID := 456
	username := "claims_test_user"

	token, err := jwt.GenerateJWT(userID, username)
	require.NoError(t, err)

	claims, err := jwt.ParseJWT(token)
	require.NoError(t, err)

	// Test all claim fields
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
	assert.Equal(t, "wb-app", claims.Issuer)
	assert.Equal(t, "456", claims.Subject)

	// Test that token is not expired
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestJWTErrorTypes(t *testing.T) {
	// Test specific error types
	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "invalid token error",
			token:   "invalid.token",
			wantErr: jwt.ErrInvalidToken,
		},
		{
			name:    "missing token error",
			token:   "",
			wantErr: jwt.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwt.ParseJWT(tt.token)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
