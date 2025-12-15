package db

import (
	"log"

	"mobile-order-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("toeic_system.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema (Could be moved to a separate migrate tool, but ok here for prototype)
	err = db.AutoMigrate(&models.SystemAdmin{}, &models.RestaurantStaff{}, &models.Restaurant{}, &models.MenuItem{})
	if err != nil {
		log.Fatal("failed to migrate schema")
	}

	return db
}
