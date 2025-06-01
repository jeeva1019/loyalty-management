package config

import (
	"fmt"
	"loyality_points/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB establishes a connection to the MySQL database using GORM.
func ConnectDB() (*gorm.DB, error) {
	// Read database configuration from TOML
	user := GetTomlStr("dbconfig", "user")
	pwd := GetTomlStr("dbconfig", "password")
	host := GetTomlStr("dbconfig", "host")
	dbName := GetTomlStr("dbconfig", "db")

	if user == "" || pwd == "" || host == "" || dbName == "" {
		return nil, fmt.Errorf("database configuration values are missing")
	}

	// Data Source Name with parseTime to correctly map time.Time fields
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pwd, host, dbName)

	// Attempt to connect to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("✅ Database connection established")

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate models: %w", err)
	}
	// Auto-migrate schema for the defined models
	fmt.Println("✅ Database migrated successfully")
	return db, nil
}

// migrate handles model migrations.
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.PointsRecord{}, &models.Transaction{})
}
