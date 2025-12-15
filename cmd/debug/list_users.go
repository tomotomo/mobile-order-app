package main

import (
	"fmt"
	"log"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/models"
)

func main() {
	database := db.Init()

	var staffs []models.RestaurantStaff
	if err := database.Find(&staffs).Error; err != nil {
		log.Fatalf("Failed to fetch staffs: %v", err)
	}

	fmt.Println("--- Restaurant Staffs ---")
	for _, s := range staffs {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Role: %s, RestaurantID: %d\n", 
			s.ID, s.Name, s.Email, s.Role, s.RestaurantID)
	}

	var admins []models.SystemAdmin
	if err := database.Find(&admins).Error; err != nil {
		log.Fatalf("Failed to fetch admins: %v", err)
	}

	fmt.Println("\n--- System Admins ---")
	for _, a := range admins {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", a.ID, a.Name, a.Email)
	}

	var restaurants []models.Restaurant
	if err := database.Find(&restaurants).Error; err != nil {
		log.Fatalf("Failed to fetch restaurants: %v", err)
	}

	fmt.Println("\n--- Restaurants ---")
	for _, r := range restaurants {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}
