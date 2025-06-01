package services

import (
	"errors"
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func (s *UserService) VerifyUser(debug *helpers.HelperStruct, userInfo models.User) (models.User, int, error) {
	debug.Info("VerifyUser(+)")
	var currentInfo models.User

	email := strings.TrimSpace(userInfo.Email)

	// Check if user with email exists
	if err := s.DB.Where("email = ?", email).First(&currentInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			debug.Error("LHVU:001", err)
			return currentInfo, http.StatusUnauthorized, errors.New("user not found or email is incorrect")
		}
		debug.Error("LHVU:DBERR", err)
		return currentInfo, http.StatusInternalServerError, errors.New("database error")
	}

	// Check password match
	if err := utils.CheckPassword(currentInfo.Password, userInfo.Password); err != nil {
		debug.Error("LHVU:002", err)
		return currentInfo, http.StatusBadRequest, errors.New("invalid password")
	}

	debug.Info("VerifyUser(-)")
	return currentInfo, 0, nil
}
