package main

import (
	"log"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/models"
)

func main() {
	log.Println("Starting schema split migration...")
	database := db.Init()

	// 1. Migrate Admin Users
	var admins []models.User
	if err := database.Where("role = ?", "admin").Find(&admins).Error; err == nil {
		for _, u := range admins {
			sa := models.SystemAdmin{
				Name:         u.Name,
				Email:        u.Email,
				PasswordHash: u.PasswordHash,
			}
			// Use FirstOrCreate to avoid dups if run multiple times
			if err := database.Where("email = ?", sa.Email).FirstOrCreate(&sa).Error; err != nil {
				log.Printf("Failed to migrate admin %s: %v", u.Email, err)
			} else {
				log.Printf("Migrated admin: %s", u.Email)
			}
		}
	}

	// 2. Migrate Restaurant Staff/Managers
	var staffs []models.User
	if err := database.Where("role IN ?", []string{"manager", "staff", "restaurant"}).Find(&staffs).Error; err == nil {
		for _, u := range staffs {
			if u.RestaurantID == nil {
				log.Printf("Skipping staff %s: No RestaurantID", u.Email)
				continue
			}

			role := models.StaffRoleMember
			if u.Role == "manager" || u.Role == "restaurant" {
				role = models.StaffRoleManager
			}

			rs := models.RestaurantStaff{
				RestaurantID: *u.RestaurantID,
				Name:         u.Name,
				Email:        u.Email,
				PasswordHash: u.PasswordHash,
				Role:         role,
			}
			if err := database.Where("email = ?", rs.Email).FirstOrCreate(&rs).Error; err != nil {
				log.Printf("Failed to migrate staff %s: %v", u.Email, err)
			} else {
				log.Printf("Migrated staff: %s", u.Email)
			}
		}
	}

	log.Println("Schema split migration complete.")
}
