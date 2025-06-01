package controllers

import (
	"encoding/json"
	"fmt"
	"loyality_points/constants"
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"
	"time"
)

func (c *UserController) LogInHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("LogInHandler(+)")

	c.SetCORS(w)

	// Handle CORS preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var (
		err      error
		userInfo models.User
	)

	// Decode the request body JSON into userInfo struct
	if err = json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		debug.Error("CLH:001", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CLH:001]-Invalid request body", nil))
		return
	}

	if err := utils.Validator(map[string]string{
		"email":    userInfo.Email,
		"password": userInfo.Password,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CLH:007]-"+err.Error(), nil))
		return
	}

	// Verify user credentials using the UserService
	userInfo, statusCode, err := c.UserService.VerifyUser(debug, userInfo)
	if err != nil {
		debug.Error("CLH:002", err)
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CLH:002]-"+err.Error(), nil))
		return
	}

	// Set access and refresh tokens if login is successful
	SetAuth(debug, userInfo, w)
	debug.Info("LogInHandler(-)")
}

func SetAuth(debug *helpers.HelperStruct, userInfo models.User, w http.ResponseWriter) {
	debug.Info("SetAuth(+)")

	// Generate access token with user's email
	accessToken, err := utils.GenerateAccessJWT(userInfo.Email)
	if err != nil {
		debug.Error("LSA:001", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[LSA:001]-"+"Failed to generate access token", nil))
		return
	}

	// Generate refresh token using user ID
	refreshToken, err := utils.GenerateRefreshJWT(fmt.Sprintf("%d", userInfo.ID))
	if err != nil {
		debug.Error("LSA:002", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[LSA:002]-"+"Failed to generate refresh token", nil))
		return
	}

	// Set the access token in the Authorization header
	w.Header().Set("Authorization", "Bearer "+accessToken)

	// Set the refresh token in an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     constants.COOKIENAME,
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(utils.RefreshTokenDuration),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "Login successful", nil))
	debug.Info("SeAuth(-)")
}
