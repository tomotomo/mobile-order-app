package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleStaff   UserRole = "staff"
	// Guests don't have a User record in this design, they are transient
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Role         UserRole `gorm:"not null"`
	RestaurantID *uint
	Restaurant   *Restaurant `gorm:"foreignKey:RestaurantID"`
}

type Restaurant struct {
	gorm.Model
	Name                 string `gorm:"not null"`
	UserID               uint   `gorm:"not null"` // Creator/Manager ID
	Users                []User `gorm:"foreignKey:RestaurantID"` // All staff including manager
	
	// Single Access Locking Mechanism
	CurrentGuestSession  string    // Session ID of the active guest
	GuestSessionExpires  time.Time // When the lock expires
}

type MenuItem struct {
	gorm.Model
	RestaurantID uint `gorm:"not null"`
	Name         string `gorm:"not null"`
	Price        int    `gorm:"not null"`
	Stock        int    `gorm:"default:0"`
	IsSoldOut    bool   `gorm:"default:false"`
}
