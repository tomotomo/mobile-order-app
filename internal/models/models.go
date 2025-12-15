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

// DEPRECATED: Use SystemAdmin or RestaurantStaff
type User struct {
	gorm.Model
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Role         UserRole
	RestaurantID *uint
}

type SystemAdmin struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}

type StaffRole string
const (
	StaffRoleManager StaffRole = "manager"
	StaffRoleMember  StaffRole = "staff"
)

type RestaurantStaff struct {
	gorm.Model
	RestaurantID uint      `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         StaffRole `gorm:"not null"` // manager or staff
}

type Restaurant struct {
	gorm.Model
	Name                 string            `gorm:"not null"`
	Staffs               []RestaurantStaff `gorm:"foreignKey:RestaurantID"` 
	
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
