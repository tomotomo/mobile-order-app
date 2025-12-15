package guest

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"mobile-order-app/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type CheckInResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"` // success, busy
}

func (h *Handler) CheckIn(c echo.Context) error {
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

func (h *Handler) GetMenu(c echo.Context) error {
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

func (h *Handler) Checkout(c echo.Context) error {
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
