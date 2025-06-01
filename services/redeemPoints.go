package services

import (
	"errors"
	"loyality_points/helpers"
	"loyality_points/models"
	"net/http"

	"gorm.io/gorm"
)

func (s *UserService) RedeemPoints(debug *helpers.HelperStruct, email string, req models.RedeemPointsRequest) (models.User, int, error) {
	debug.Info("RedeemPoints(+)")

	// 1. Retrieve user record by email
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		debug.Error("RPHRP:001", err)
		return user, http.StatusNotFound, errors.New("user not found")
	}

	// 2. Check if user has enough points to redeem
	if user.PointsBalance < req.PointsToRedeem {
		debug.Error("RPHRP:002", "insufficient points")
		return user, http.StatusBadRequest, errors.New("insufficient points")
	}

	// 3. Create a new points record to deduct points (redeem)
	pointsRecord := models.PointsRecord{
		UserID: user.ID,
		Points: -req.PointsToRedeem,
		Type:   "redeem",
		Reason: req.Reason,
	}

	// 4. Update user's points balance by subtracting redeemed points atomically
	if err := s.DB.Create(&pointsRecord).Error; err != nil {
		debug.Error("RPHRP:003", err.Error())
		return user, http.StatusInternalServerError, err
	}

	// 5. Fetch the updated user record to return fresh data
	if err := s.DB.
		Model(&models.User{}).
		Where("id = ?", user.ID).
		UpdateColumn("points_balance", gorm.Expr("points_balance - ?", req.PointsToRedeem)).Error; err != nil {
		debug.Error("RPHRP:004", err.Error())
		return user, http.StatusInternalServerError, err
	}

	if err := s.DB.First(&user, user.ID).Error; err != nil {
		debug.Error("RPHRP:005", err.Error())
		return user, http.StatusInternalServerError, err
	}

	debug.Info("RedeemPoints(-)")
	return user, http.StatusOK, nil
}
