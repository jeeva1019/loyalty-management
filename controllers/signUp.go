package controllers

import (
	"encoding/json"
	"fmt"
	"loyality_points/constants"
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"
)

func (c *UserController) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("SignUp(+)")

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

	// Decode JSON request body into userInfo struct
	if err = json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		debug.Error("CSU:001", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CSU:001]-"+err.Error(), nil))
		return
	}

	if err := utils.Validator(map[string]string{
		"email":    userInfo.Email,
		"password": userInfo.Password,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CSU:007]-"+err.Error(), nil))
		return
	}

	// Call UserService to add a new user
	if statusCode, err := c.UserService.AddUser(debug, userInfo); err != nil {
		debug.Error("CSU:002", err)
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CSU:002]-"+err.Error(), nil))
		return
	}

	// Respond with 202 Accepted and success message on successful user creation
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "User Created Successfully", nil))
	debug.Info("SignUp(-)")
}
