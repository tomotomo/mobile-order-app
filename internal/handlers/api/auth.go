package api

import (
	"net/http"
	"time"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Try System Admin
	var admin models.SystemAdmin
	if err := h.DB.Where("email = ?", req.Email).First(&admin).Error; err == nil {
		if auth.CheckPasswordHash(req.Password, admin.PasswordHash) {
			token, _ := auth.GenerateToken(admin.ID, "admin")
			// Update LastLoginAt
			now := time.Now()
			h.DB.Model(&admin).Update("last_login_at", now)
			return c.JSON(http.StatusOK, map[string]string{"token": token, "role": "admin"})
		}
	}

	// Try Restaurant Staff
	var staff models.RestaurantStaff
	if err := h.DB.Where("email = ?", req.Email).First(&staff).Error; err == nil {
		if auth.CheckPasswordHash(req.Password, staff.PasswordHash) {
			token, _ := auth.GenerateToken(staff.ID, string(staff.Role))
			// Update LastLoginAt
			now := time.Now()
			h.DB.Model(&staff).Update("last_login_at", now)
			return c.JSON(http.StatusOK, map[string]string{"token": token, "role": string(staff.Role)})
		}
	}

	return c.JSON(http.StatusUnauthorized, "invalid credentials")
}
