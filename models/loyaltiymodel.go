package models

import "time"

type User struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Email         string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password      string `gorm:"type:varchar(255);not null"`
	PointsBalance int64
	Transactions  []Transaction  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PointsRecords []PointsRecord `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Transaction struct {
	ID                uint      `gorm:"primaryKey"`
	TransactionID     string    `gorm:"type:varchar(100);uniqueIndex" json:"transaction_id"`
	UserID            uint      `json:"user_id"`
	User              User      `gorm:"foreignKey:UserID"`
	TransactionAmount float64   `json:"transaction_amount"`
	Category          string    `gorm:"type:varchar(100)" json:"category"`
	TransactionDate   time.Time `json:"transaction_date"`
	ProductCode       string    `gorm:"type:varchar(100)" json:"product_code"`
	CreatedAt         time.Time
	PointsRecords     []PointsRecord `gorm:"foreignKey:RelatedTransactionID"`
}

type PointsRecord struct {
	ID                   uint `gorm:"primaryKey"`
	UserID               uint
	User                 User `gorm:"foreignKey:UserID"`
	Points               int64
	Type                 string `gorm:"type:varchar(50)"`
	Reason               string `gorm:"type:varchar(255)"`
	CreatedAt            time.Time
	RelatedTransactionID *uint
	RelatedTransaction   *Transaction `gorm:"foreignKey:RelatedTransactionID"`
}

type RedeemPointsRequest struct {
	PointsToRedeem int64  `json:"points_to_redeem"`
	Reason         string `json:"reason"`
}

type PointsHistoryDTO struct {
	Email     string    `json:"email"`
	Points    int64     `json:"points"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Category  string    `json:"category,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
