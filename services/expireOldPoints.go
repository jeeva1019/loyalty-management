package services

import (
	"loyality_points/helpers"
	"loyality_points/models"
	"time"

	"gorm.io/gorm"
)

func (s *UserService) ExpireOldPoints(debug *helpers.HelperStruct) {
	debug.Info("expireOldPoints(+)")

	// Define cutoff date: points older than this will expire (1 year ago)
	cutoffDate := time.Now().AddDate(-1, 0, 0)

	// Find all "earn" points records older than cutoffDate with positive points
	var oldPoints []models.PointsRecord
	if err := s.DB.Where("type = ? AND created_at <= ? AND points > 0", "earn", cutoffDate).Find(&oldPoints).Error; err != nil {
		debug.Error("SEOP:001", "Error finding old points to expire:", err)
		return
	}

	// Loop through each points record to expire them individually
	for _, pr := range oldPoints {
		tx := s.DB.Begin()

		expireRecord := models.PointsRecord{
			UserID: pr.UserID,
			Points: -pr.Points,
			Type:   "expire",
			Reason: "expired",
		}

		if err := tx.Create(&expireRecord).Error; err != nil {
			tx.Rollback()
			debug.Error("SEOP:002", "Error creating expire record:", err)
			continue
		}

		if err := tx.Model(&models.User{}).Where("id = ?", pr.UserID).UpdateColumn("points_balance", gorm.Expr("points_balance - ?", pr.Points)).Error; err != nil {
			tx.Rollback()
			debug.Error("SEOP:003", "Error updating user balance:", err)
			continue
		}

		tx.Commit()
	}
	debug.Info("expireOldPoints(-)")
}
