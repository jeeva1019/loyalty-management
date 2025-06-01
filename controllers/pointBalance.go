package controllers

import (
	"fmt"
	"loyality_points/constants"
	"loyality_points/helpers"
	"loyality_points/utils"
	"net/http"
)

func (c *UserController) PointsBalanceHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("PointsBalanceHandler(+)")

	// Extract email from JWT token stored in request context
	email, ok := r.Context().Value(constants.TOKENKEY).(string)
	if !ok || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:001]-Unauthorized", nil))
		return
	}

	// Extract pagination parameters from query string
	start := r.URL.Query().Get("page")
	end := r.URL.Query().Get("page_size")

	if err := utils.Validator(map[string]string{
		"page":      start,
		"page_size": end,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CSU:007]-"+err.Error(), nil))
		return
	}

	// Fetch the user's point balance using UserService
	pointResp, statusCode, err := c.UserService.GetPointBalance(debug, email, start, end)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:002]-"+err.Error(), nil))
		return
	}

	// Send successful response with points data
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "", pointResp))
	debug.Info("PointsBalanceHandler(-)")
}
