package restaurant_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/handlers/restaurant"
	"mobile-order-app/internal/models"
	"mobile-order-app/internal/testutils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockEmailSender
type MockEmailSender struct {
	SentEmails []struct {
		To, Subject, Body string
	}
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	m.SentEmails = append(m.SentEmails, struct{ To, Subject, Body string }{to, subject, body})
	return nil
}

func TestInviteStaff(t *testing.T) {
	db := testutils.SetupTestDB()
	mockEmail := &MockEmailSender{}
	h := restaurant.NewHandler(db, mockEmail)
	e := echo.New()

	// Seed Manager and Restaurant
	passHash, _ := auth.HashPassword("pass")
	rest := models.Restaurant{Name: "Test Rest"}
	db.Create(&rest)

	manager := models.RestaurantStaff{
		RestaurantID: rest.ID,
		Name:         "Manager",
		Email:        "manager@test.com",
		PasswordHash: passHash,
		Role:         models.StaffRoleManager,
	}
	db.Create(&manager)

	// Seed Staff (for negative test)
	staff := models.RestaurantStaff{
		RestaurantID: rest.ID,
		Name:         "Staff",
		Email:        "staff@test.com",
		PasswordHash: passHash,
		Role:         models.StaffRoleMember,
	}
	db.Create(&staff)

	// Setup JWT Middleware for Context
	// Helper to create authenticated context
	createContext := func(userID uint, role string, req *http.Request) echo.Context {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Manually set JWT token (simulating middleware success)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.JWTClaims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{},
		})
		c.Set("user", token)
		return c
	}

	t.Run("Manager Can Invite Staff", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"email": "newuser@test.com",
			"name":  "New User",
		})
		req := httptest.NewRequest(http.MethodPost, "/staff/invite", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		
		c := createContext(manager.ID, string(models.StaffRoleManager), req)

		if assert.NoError(t, h.InviteStaff(c)) {
			assert.Equal(t, http.StatusCreated, c.Response().Status)
			
			// Verify DB
			var newUser models.RestaurantStaff
			db.Where("email = ?", "newuser@test.com").First(&newUser)
			assert.Equal(t, "New User", newUser.Name)
			assert.Equal(t, models.StaffRoleMember, newUser.Role) // Created as Staff
			assert.Equal(t, rest.ID, newUser.RestaurantID)

			// Verify Email
			assert.Len(t, mockEmail.SentEmails, 1)
			assert.Equal(t, "newuser@test.com", mockEmail.SentEmails[0].To)
		}
	})

	t.Run("Staff Cannot Invite", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"email": "shouldfail@test.com",
			"name":  "Fail User",
		})
		req := httptest.NewRequest(http.MethodPost, "/staff/invite", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		
		c := createContext(staff.ID, string(models.StaffRoleMember), req)

		// Handler returns 403 error inside JSON or plain
		// Check implementation: return c.JSON(403, ...)
		h.InviteStaff(c)
		assert.Equal(t, http.StatusForbidden, c.Response().Status)
	})
}
