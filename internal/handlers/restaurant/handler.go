package restaurant

import (
	"net/http"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
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
