package controllers

import (
	"fmt"
	"loyality_points/constants"
	"loyality_points/helpers"
	"loyality_points/utils"
	"net/http"
)

func (c *UserController) PointsHistoryHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("PointsHistoryHandler(+)")

	// Extract email from JWT token stored in request context
	email, ok := r.Context().Value(constants.TOKENKEY).(string)
	if !ok || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:001]-Unauthorized", nil))
		return
	}

	// Parse query parameters for filtering the points history
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	txType := r.URL.Query().Get("txtype")

	if err := utils.Validator(map[string]string{
		"start_date": startDate,
		"end_date":   endDate,
		"start":      start,
		"end":        end,
		"txtype":     txType,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:007]-"+err.Error(), nil))
		return
	}

	// Call service to get points history based on filters
	resp, statusCode, err := c.UserService.GetPointsHistory(debug, email, txType, startDate, endDate, start, end)
	if err != nil {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:002]-"+err.Error(), nil))
		return
	}

	// Write success response with the points history data
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "", resp))
	debug.Info("PointsHistoryHandler(-)")
}
