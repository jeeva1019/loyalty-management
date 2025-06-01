package utils

import (
	"errors"
	"loyality_points/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration  = time.Hour          // 1 hour
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
)

// GenerateAccessJWT creates a short-lived token using email
func GenerateAccessJWT(email string) (string, error) {
	var jwtKey = []byte(config.GetTomlStr("common", "secret"))
	claims := jwt.MapClaims{
		"email": email,
		"type":  "access",
		"exp":   time.Now().Add(AccessTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// GenerateRefreshJWT creates a long-lived refresh token using userID
func GenerateRefreshJWT(userID string) (string, error) {
	var jwtKey = []byte(config.GetTomlStr("common", "secret"))
	claims := jwt.MapClaims{
		"userID": userID,
		"type":   "refresh",
		"exp":    time.Now().Add(RefreshTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT validates a token and returns claim type and identifier
func ValidateJWT(tokenStr string) (map[string]string, error) {
	var jwtKey = []byte(config.GetTomlStr("common", "secret"))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	typ, ok := claims["type"].(string)
	if !ok {
		return nil, errors.New("missing token type")
	}

	result := map[string]string{"type": typ}

	switch typ {
	case "access":
		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("missing email in access token")
		}
		result["email"] = email
	case "refresh":
		userID, ok := claims["userID"].(string)
		if !ok {
			return nil, errors.New("missing userID in refresh token")
		}
		result["userID"] = userID
	default:
		return nil, errors.New("unknown token type")
	}

	return result, nil
}
