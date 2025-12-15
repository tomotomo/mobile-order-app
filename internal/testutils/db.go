package testutils

import (
	"mobile-order-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite DB and migrates schema
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate schema
	err = db.AutoMigrate(
		&models.SystemAdmin{},
		&models.RestaurantStaff{},
		&models.Restaurant{},
		&models.MenuItem{},
	)
	if err != nil {
		panic("failed to migrate schema")
	}

	return db
}
