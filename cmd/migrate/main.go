package main

import (
	"log"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/models"
)

func main() {
	log.Println("Starting migration...")
	database := db.Init()

	// Update Users with Role "restaurant" to "manager"
	// Old schema used "restaurant" string
	result := database.Model(&models.User{}).Where("role = ?", "restaurant").Update("role", models.RoleManager)

	if result.Error != nil {
		log.Fatalf("Migration failed: %v", result.Error)
	}

	log.Printf("Migration complete. Updated %d users from 'restaurant' to 'manager'.", result.RowsAffected)
}
