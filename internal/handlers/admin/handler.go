package admin

import (
	"fmt"
	"net/http"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/email"
	"mobile-order-app/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	Email email.Sender
}

func NewHandler(db *gorm.DB, emailSender email.Sender) *Handler {
	return &Handler{
		DB:    db,
		Email: emailSender,
	}
}

func (h *Handler) Dashboard(c echo.Context) error {
	var rests []models.Restaurant
	if err := h.DB.Preload("User").Find(&rests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, rests)
}

type CreateRestRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) CreateRestaurant(c echo.Context) error {
	req := new(CreateRestRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	hash, _ := auth.HashPassword(req.Password)
	tx := h.DB.Begin()

	// Create Manager (RestaurantStaff)
	manager := models.RestaurantStaff{
		Name:         req.Name + " Manager",
		Email:        req.Email,
		PasswordHash: hash,
		Role:         models.StaffRoleManager,
	}
	// Create Restaurant first
	rest := models.Restaurant{
		Name: req.Name,
	}
	if err := tx.Create(&rest).Error; err != nil {
		tx.Rollback(); return c.JSON(500, err)
	}

	manager.RestaurantID = rest.ID
	if err := tx.Create(&manager).Error; err != nil {
		tx.Rollback(); return c.JSON(500, err)
	}

	tx.Commit()

	// Send Welcome Email
	subject := "Welcome to Mobile Order App"
	body := fmt.Sprintf("Hello %s,\n\nYour restaurant '%s' has been registered.\nLogin with:\nEmail: %s\nPassword: %s", 
		manager.Name, rest.Name, manager.Email, req.Password)
	
	go h.Email.Send(manager.Email, subject, body)

	return c.JSON(http.StatusCreated, rest)
}

func (h *Handler) SuspendRestaurant(c echo.Context) error {
	id := c.Param("id")
	h.DB.Delete(&models.Restaurant{}, id)
	return c.JSON(http.StatusOK, "suspended")
}
