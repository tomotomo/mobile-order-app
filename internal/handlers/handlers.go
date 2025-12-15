package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

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

// --- Auth ---

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (h *Handler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
		"role":  string(user.Role),
	})
}

// --- AdminHandlers ---

func (h *Handler) AdminDashboard(c echo.Context) error {
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

	user := models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name + " Manager",
		Role:         models.RoleRestaurant,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback(); return c.JSON(500, err)
	}

	rest := models.Restaurant{
		Name: req.Name,
		UserID: user.ID,
	}
	if err := tx.Create(&rest).Error; err != nil {
		tx.Rollback(); return c.JSON(500, err)
	}

	user.RestaurantID = &rest.ID
	tx.Save(&user)
	tx.Commit()

	return c.JSON(http.StatusCreated, rest)
}

func (h *Handler) SuspendRestaurant(c echo.Context) error {
	id := c.Param("id")
	// Simplification: Delete
	h.DB.Delete(&models.Restaurant{}, id)
	return c.JSON(http.StatusOK, "suspended")
}

// --- Restaurant Handlers ---

func (h *Handler) RestaurantDashboard(c echo.Context) error {
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

// --- Guest Handlers (The Core Logic) ---

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type CheckInResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"` // success, busy
}

func (h *Handler) GuestCheckIn(c echo.Context) error {
	restID := c.Param("id")
	var rest models.Restaurant
	if err := h.DB.First(&rest, restID).Error; err != nil {
		return c.JSON(404, "restaurant not found")
	}

	// Logic: Check if anyone is concurrently active
	now := time.Now()
	isLocked := rest.CurrentGuestSession != "" && rest.GuestSessionExpires.After(now)

	if isLocked {
		return c.JSON(423, map[string]string{"status": "busy", "message": "Someone else is ordering right now."})
	}

	// Acquire Lock
	sessionID := generateSessionID()
	rest.CurrentGuestSession = sessionID
	rest.GuestSessionExpires = now.Add(5 * time.Minute) // 5 min timeout for prototype
	h.DB.Save(&rest)

	return c.JSON(200, CheckInResponse{SessionID: sessionID, Status: "success"})
}

func (h *Handler) GuestGetMenu(c echo.Context) error {
	restID := c.Param("id")
	sessionID := c.QueryParam("session_id")

	var rest models.Restaurant
	h.DB.First(&rest, restID)

	// Validate Lock
	if rest.CurrentGuestSession != sessionID || rest.GuestSessionExpires.Before(time.Now()) {
		return c.JSON(401, "session expired or invalid")
	}

	// Refresh Lock? (Optional)
	rest.GuestSessionExpires = time.Now().Add(5 * time.Minute)
	h.DB.Save(&rest)

	var items []models.MenuItem
	h.DB.Where("restaurant_id = ?", restID).Find(&items)
	return c.JSON(200, items)
}

func (h *Handler) GuestCheckout(c echo.Context) error {
	restID := c.Param("id")
	sessionID := c.QueryParam("session_id") // or body

	var rest models.Restaurant
	h.DB.First(&rest, restID)

	if rest.CurrentGuestSession == sessionID {
		// Release Lock
		rest.CurrentGuestSession = ""
		rest.GuestSessionExpires = time.Time{} // zero
		h.DB.Save(&rest)
		return c.JSON(200, "checked out")
	}
	return c.JSON(400, "invalid session")
}
