package main

import (
	"log"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/models"
)

func main() {
	log.Println("Starting migration: Add LastLoginAt...")
	database := db.Init()

	// GORM AutoMigrate will add missing columns
	if err := database.AutoMigrate(&models.SystemAdmin{}, &models.RestaurantStaff{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration complete.")
}
