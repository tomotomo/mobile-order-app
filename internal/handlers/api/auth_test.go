package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/handlers/api"
	"mobile-order-app/internal/models"
	"mobile-order-app/internal/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	db := testutils.SetupTestDB()
	h := api.NewAuthHandler(db)
	e := echo.New()

	// Seed Admin
	passHash, _ := auth.HashPassword("admin123")
	admin := models.SystemAdmin{
		Name:         "AdminUser",
		Email:        "admin@test.com",
		PasswordHash: passHash,
	}
	db.Create(&admin)

	// Seed Staff
	rest := models.Restaurant{Name: "Test Rest"}
	db.Create(&rest)
	staff := models.RestaurantStaff{
		RestaurantID: rest.ID,
		Name:         "StaffUser",
		Email:        "staff@test.com",
		PasswordHash: passHash,
		Role:         models.StaffRoleMember,
	}
	db.Create(&staff)

	t.Run("Admin Login Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"email":    "admin@test.com",
			"password": "admin123",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var res map[string]string
			json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NotEmpty(t, res["token"])
			assert.Equal(t, "admin", res["role"])
		}
	})

	t.Run("Staff Login Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"email":    "staff@test.com",
			"password": "admin123",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var res map[string]string
			json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Equal(t, "staff", res["role"])
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"email":    "admin@test.com",
			"password": "wrongpassword",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})
}
