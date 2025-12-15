package restaurant

import (
	"fmt"
	"net/http"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/email"
	"mobile-order-app/internal/models"

	"github.com/golang-jwt/jwt/v5"
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
	userTok := c.Get("user").(*jwt.Token)
	claims := userTok.Claims.(*auth.JWTClaims)
	
	var user models.User
	h.DB.First(&user, claims.UserID)
	if user.RestaurantID == nil { return c.JSON(403, "no restaurant") }

	var items []models.MenuItem
	h.DB.Where("restaurant_id = ?", *user.RestaurantID).Find(&items)
	return c.JSON(http.StatusOK, items)
}

type MenuItemReq struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func (h *Handler) CreateMenuItem(c echo.Context) error {
	userTok := c.Get("user").(*jwt.Token)
	claims := userTok.Claims.(*auth.JWTClaims)
	var user models.User
	h.DB.First(&user, claims.UserID)

	req := new(MenuItemReq)
	c.Bind(req)

	item := models.MenuItem{
		RestaurantID: *user.RestaurantID,
		Name:         req.Name,
		Price:        req.Price,
		Stock:        req.Stock,
	}
	h.DB.Create(&item)
	return c.JSON(201, item)
}

type InviteRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) InviteStaff(c echo.Context) error {
	userTok := c.Get("user").(*jwt.Token)
	claims := userTok.Claims.(*auth.JWTClaims)

	// Authorization Check: Only Managers can invite
	if claims.Role != string(models.RoleManager) {
		return c.JSON(403, "only managers can invite staff")
	}

	var manager models.User
	h.DB.First(&manager, claims.UserID)

	req := new(InviteRequest)
	if err := c.Bind(req); err != nil { return c.JSON(400, err) }

	// Generate temp password (random 8 chars for MVP)
	tempPass := "staff123" // TODO: Randomize
	hash, _ := auth.HashPassword(tempPass)

	staff := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		Role:         models.RoleStaff,
		RestaurantID: manager.RestaurantID,
	}

	if err := h.DB.Create(&staff).Error; err != nil {
		return c.JSON(500, err)
	}

	// Send Email
	subject := "Invitation to Mobile Order App"
	body := fmt.Sprintf("Hello %s,\n\nYou have been invited to join '%s' as staff.\nLogin with:\nEmail: %s\nPassword: %s",
		req.Name, "Your Restaurant", req.Email, tempPass) // Ideally fetch Restaurant Name

	go h.Email.Send(req.Email, subject, body)

	return c.JSON(201, staff)
}

func (h *Handler) UpdateStock(c echo.Context) error {
	id := c.Param("id")
	var item models.MenuItem
	if err := h.DB.First(&item, id).Error; err != nil { return c.JSON(404, "not found") }
	
	type StockReq struct { Stock int `json:"stock"`; IsSoldOut bool `json:"is_sold_out"` }
	req := new(StockReq)
	c.Bind(req)

	item.Stock = req.Stock
	item.IsSoldOut = req.IsSoldOut
	h.DB.Save(&item)
	return c.JSON(200, item)
}
