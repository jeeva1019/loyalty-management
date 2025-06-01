package services

import (
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"

	"gorm.io/gorm"
)

func (s *UserService) GetPointBalance(debug *helpers.HelperStruct, email, start, end string) (map[string]any, int, error) {
	debug.Info("GetPointBalance(+)")

	respMap := make(map[string]any)

	// 1. Get user by email
	var user models.User
	if err := s.DB.Model(&models.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		debug.Error("PBHGPB:001", err)
		return nil, http.StatusInternalServerError, err
	}

	// 2. Parse pagination parameters
	page, pageSize, offset := utils.GetPaginationValue(start, end)

	// 3. Fetch paginated points history
	var pointsRecords []models.PointsRecord
	if err := s.DB.
		Where("user_id = ?", user.ID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Preload("RelatedTransaction", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category")
		}).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&pointsRecords).Error; err != nil {
		debug.Error("PBHGPB:002", err)
		return nil, http.StatusInternalServerError, err
	}

	// 4. Format the response
	var history = responseFormat(pointsRecords)

	// 5. Build response
	respMap["points_balance"] = user.PointsBalance
	respMap["points_history"] = history
	respMap["page"] = page
	respMap["page_size"] = pageSize

	debug.Info("GetPointBalance(-)")
	return respMap, 0, nil
}
