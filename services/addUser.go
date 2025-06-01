package services

import (
	"errors"
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"

	"gorm.io/gorm"
)

func (s *UserService) AddUser(debug *helpers.HelperStruct, userInfo models.User) (int, error) {
	debug.Info("AddUser(+)")

	// Validate email
	if err := utils.ValidateEmail(userInfo.Email); err != nil {
		debug.Error("SHAU:001", err)
		return http.StatusBadRequest, err
	}

	// Validate password
	if err := utils.ValidatePassword(userInfo.Password); err != nil {
		debug.Error("SHAU:002", err)
		return http.StatusBadRequest, err
	}

	// Hash password
	hashedPwd, err := utils.HashPassword(userInfo.Password)
	if err != nil {
		debug.Error("SHAU:003", err)
		return http.StatusInternalServerError, err
	}
	userInfo.Password = hashedPwd

	// Check if user already exists
	var existing models.User
	err = s.DB.Where("email = ?", userInfo.Email).First(&existing).Error
	if err == nil {
		debug.Error("SHAU:005", errors.New("user already exists"))
		return http.StatusConflict, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		debug.Error("SHAU:006", err)
		return http.StatusInternalServerError, err
	}

	// Create new user
	if err := s.DB.Create(&userInfo).Error; err != nil {
		debug.Error("SHAU:004", err)
		return http.StatusInternalServerError, err
	}

	debug.Info("AddUser(-)")
	return http.StatusCreated, nil
}
