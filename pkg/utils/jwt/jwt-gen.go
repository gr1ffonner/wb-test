package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTSecret       = []byte("TheSecretKeyAlwayssafe")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
	ErrMissingToken = errors.New("missing token")
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token with user information
func GenerateJWT(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wb-app",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseJWT validates and parses a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	// Check if it's Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidToken
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	user, ok := ctx.Value("user").(*Claims)
	return user, ok
}

// RefreshToken generates a new token with extended expiration
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	// Calculate new expiration: extend from original expiration time
	var newExpiration time.Time
	if claims.ExpiresAt != nil {
		// Extend by 24 hours from the original expiration
		newExpiration = claims.ExpiresAt.Time.Add(24 * time.Hour)
	} else {
		// Fallback: 24 hours from now
		newExpiration = time.Now().Add(24 * time.Hour)
	}

	// Create new claims with extended expiration
	newClaims := Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(newExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wb-app",
			Subject:   fmt.Sprintf("%d", claims.UserID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString(JWTSecret)
}

// ValidateToken validates token without parsing claims
func ValidateToken(tokenString string) error {
	_, err := ParseJWT(tokenString)
	return err
}
