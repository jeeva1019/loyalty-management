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

func (c *UserController) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	debug := new(helpers.HelperStruct)
	debug.Init()
	debug.Info("TransactionHandler(+)")

	// Extract the user's email from the request context (JWT token)
	email, ok := r.Context().Value(constants.TOKENKEY).(string)
	if !ok || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CTH:001]-Unauthorized", nil))
		return
	}

	var (
		transactionInfo models.Transaction
		err             error
	)

	// Decode the JSON body into transactionInfo struct
	if err = json.NewDecoder(r.Body).Decode(&transactionInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CTH:002]-Invalid request body", nil))
		return
	}

	if err := utils.Validator(map[string]string{
		"transaction_id":     transactionInfo.TransactionID,
		"transaction_amount": fmt.Sprintf("%v", transactionInfo.TransactionAmount),
		"category":           transactionInfo.Category,
		"transaction_date":   transactionInfo.TransactionDate.String(),
		"product_code":       transactionInfo.ProductCode,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CPBH:007]-"+err.Error(), nil))
		return
	}

	// Call UserService to add the transaction and get points earned
	pointsEarned, err := c.UserService.AddTransaction(debug, transactionInfo, email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, utils.ResponseConstructor(constants.ERROR, "[CTH:002]-"+err.Error(), nil))
		return
	}

	// Send success response with points earned from this transaction
	fmt.Fprintln(w, utils.ResponseConstructor(constants.SUCCESS, "transaction recorded", map[string]any{
		"points_earned": pointsEarned,
	}))

	debug.Info("TransactionHandler(-)")
}
