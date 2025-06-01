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

func (c *UserController) RedeemPointsHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("RedeemPointsHandler(+)")

	// Extract user email from JWT token stored in request context
	email, ok := r.Context().Value(constants.TOKENKEY).(string)
	if !ok || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CRPH:001]-Unauthorized", nil))
		return
	}

	// Parse the JSON request body into RedeemPointsRequest struct
	var redeemReq models.RedeemPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&redeemReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CRPH:002]-Invalid request body", nil))
		return
	}

	if err := utils.Validator(map[string]string{
		"points_to_redeem": fmt.Sprintf("%d", redeemReq.PointsToRedeem),
		"reason":           redeemReq.Reason,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:007]-"+err.Error(), nil))
		return
	}

	// Call the UserService to redeem points for the user
	user, statusCode, err := c.UserService.RedeemPoints(debug, email, redeemReq)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CRPH:003]-"+err.Error(), nil))
		return
	}

	// Send success response with updated points balance
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "points redeemed", map[string]any{
		"remaining_balance": user.PointsBalance,
	}))
	debug.Info("RedeemPointsHandler(-)")
}
