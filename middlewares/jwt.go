package middleware

import (
	"context"
	"errors"
	"fmt"
	"loyality_points/constants"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"
	"strings"
)

func GetCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("no session cookie found")
		}
		return "", errors.New("error reading cookie")
	}
	return cookie.Value, nil
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Use specific origin in production
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
}

func (m Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)

		// Handle preflight OPTIONS requests (skip auth logic)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 1. Get Refresh Token from cookie
		refreshStr, err := GetCookie(r, constants.COOKIENAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "No valid refresh token provided", nil))
			return
		}

		refreshClaims, err := utils.ValidateJWT(refreshStr)
		if err != nil || refreshClaims["type"] != "refresh" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "Invalid or expired refresh token", nil))
			return
		}

		// 2. Get Access Token from Authorization header
		authHeader := r.Header.Get("Authorization")
		var accessClaims map[string]string
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken := strings.TrimPrefix(authHeader, "Bearer ")
			accessClaims, err = utils.ValidateJWT(accessToken)
		} else {
			// No access token provided
			err = errors.New("no access token provided")
		}

		// 3. If Access Token is valid, continue request with user email in context
		if err == nil && accessClaims["type"] == "access" {
			ctx := context.WithValue(r.Context(), constants.TOKENKEY, accessClaims["email"])
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 4. If Access Token is invalid or expired but Refresh Token is valid
		userID := refreshClaims["userID"]
		var user models.User
		if err := m.DB.First(&user, "id = ?", userID).Error; err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "User not found", nil))
			return
		}

		// Generate new Access Token
		newAccessToken, err := utils.GenerateAccessJWT(user.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "Failed to create new access token", nil))
			return
		}

		// Set new access token in response header for client to use next time
		w.Header().Set("X-New-Access-Token", newAccessToken)

		// Continue request with user email in context
		ctx := context.WithValue(r.Context(), constants.TOKENKEY, user.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
