package services

import (
	"loyality_points/config"
	"loyality_points/helpers"
	"loyality_points/models"
	"strconv"

	"gorm.io/gorm"
)

func (s *UserService) AddTransaction(debug *helpers.HelperStruct, transactionInfo models.Transaction, email string) (int64, error) {
	debug.Info("AddTransaction(+)")

	// 1. Retrieve user record by email from the database
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		debug.Error("THAT:000", err)
		return 0, err
	}

	// 2. Load category from configuration 
	categoryList := config.GetTomlMap("loyalityponits")
	categoryPointStr, ok := categoryList[transactionInfo.Category]
	if !ok {
		categoryPointStr = "1.0"
	}

	points, err := strconv.ParseFloat(categoryPointStr, 64)
	if err != nil {
		debug.Error("THAT:001", err)
		return 0, err
	}

	pointsEarned := int64((transactionInfo.TransactionAmount * points) / 2)

	// 3. Create transaction with resolved user ID
	transaction := models.Transaction{
		TransactionID:     transactionInfo.TransactionID,
		UserID:            user.ID,
		TransactionAmount: transactionInfo.TransactionAmount,
		Category:          transactionInfo.Category,
		TransactionDate:   transactionInfo.TransactionDate,
		ProductCode:       transactionInfo.ProductCode,
	}

	if err := s.DB.Create(&transaction).Error; err != nil {
		debug.Error("THAT:002", err)
		return 0, err
	}

	// 4. Create points record
	pointsRecord := models.PointsRecord{
		UserID:               user.ID,
		Points:               pointsEarned,
		Type:                 "earn",
		Reason:               transactionInfo.Category,
		RelatedTransactionID: &transaction.ID,
	}

	if err := s.DB.Create(&pointsRecord).Error; err != nil {
		debug.Error("THAT:003", err)
		return 0, err
	}

	// 5. Update user's points balance
	if err := s.DB.
		Model(&models.User{}).
		Where("id = ?", user.ID).
		UpdateColumn("points_balance", gorm.Expr("points_balance + ?", pointsEarned)).Error; err != nil {
		debug.Error("THAT:004", err)
		return 0, err
	}

	debug.Info("AddTransaction(-)")
	return pointsEarned, nil
}
