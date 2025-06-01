package services

import (
	"errors"
	"loyality_points/helpers"
	"loyality_points/models"
	"loyality_points/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func (s *UserService) GetPointsHistory(debug *helpers.HelperStruct, email, txType, startDateStr, endDateStr, start, end string) (map[string]any, int, error) {
	var (
		userInfo models.User
		resp     = make(map[string]any)
		err      error
	)

	// 1. Retrieve the user by email from the database
	if err = s.DB.Where("email = ?", email).First(&userInfo).Error; err != nil {
		debug.Error("PHGPH:001", err)
		return resp, http.StatusInternalServerError, err
	}

	var startDate, endDate time.Time

	// 2. Parse startDateStr if provided, return error if invalid format
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			debug.Error("PHGPH:002", err)
			return resp, http.StatusBadRequest, errors.New("invalid start_date format")
		}
	}

	// 3. Parse endDateStr if provided, return error if invalid format
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			debug.Error("PHGPH:003", err)
			return resp, http.StatusBadRequest, errors.New("invalid end_date format")
		}
	}

	// 4. Extract pagination parameters (page, pageSize, offset) using helper util
	page, pageSize, offset := utils.GetPaginationValue(start, end)

	// 5. Build base query filtering by user ID
	query := s.DB.Where("user_id = ?", userInfo.ID)

	// 6. If transaction type filter is provided, validate and apply it
	if txType != "" {
		validTypes := map[string]bool{"earn": true, "redeem": true, "expire": true, "all": true}
		if !validTypes[txType] {
			return resp, http.StatusBadRequest, errors.New("invalid transaction type")
		}
		if txType != "all" {
			query = query.Where("type = ?", txType)
		}
	}
	if !startDate.IsZero() {
		query = query.Where("date(created_at) >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date(created_at) <= ?", endDate)
	}

	// 7. Execute the query with eager loading related User and Transaction,
	//    order by created_at descending, apply pagination with Offset and Limit
	var history []models.PointsRecord
	if err := query.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Preload("RelatedTransaction", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category")
		}).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&history).Error; err != nil {
		debug.Error("PHGPH:004", err)
		return resp, http.StatusInternalServerError, errors.New("failed to fetch points history")
	}

	// 8. Format raw data into DTOs for the response
	var dtoList = responseFormat(history)

	// 9. Get total count of matching records (without pagination)
	var total int64
	_ = query.Model(&models.PointsRecord{}).Count(&total)

	resp["points_history"] = dtoList
	resp["page"] = page
	resp["page_size"] = pageSize
	resp["total_records"] = total

	return resp, http.StatusOK, nil
}

func responseFormat(history []models.PointsRecord) []models.PointsHistoryDTO {
	var dtoList []models.PointsHistoryDTO
	for _, record := range history {
		dto := models.PointsHistoryDTO{
			Email:     record.User.Email,
			Points:    record.Points,
			Type:      record.Type,
			Reason:    record.Reason,
			CreatedAt: record.CreatedAt,
		}

		if record.RelatedTransaction != nil {
			dto.Category = record.RelatedTransaction.Category
		}

		dtoList = append(dtoList, dto)
	}
	return dtoList
}
