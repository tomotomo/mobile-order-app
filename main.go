package main

import (
	"log"

	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/handlers"
	"mobile-order-app/internal/models"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Database setup
	db, err := gorm.Open(sqlite.Open("toeic_system.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Restaurant{}, &models.MenuItem{})
	if err != nil {
		log.Fatal("failed to migrate schema")
	}

	// Seed Admin User
	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)
	if adminCount == 0 {
		hash, _ := auth.HashPassword("admin123")
		admin := models.User{
			Email:        "admin@example.com",
			PasswordHash: hash,
			Name:         "System Admin",
			Role:         models.RoleAdmin,
		}
		db.Create(&admin)
		log.Println("Admin user created: admin@example.com / admin123")
	}

	// Handlers
	h := handlers.NewHandler(db)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static
	e.Static("/", "views")
	
	// Routes
	// e.GET("/", h.HealthCheck) // Replaced by static index.html
	e.POST("/login", h.Login)

	// Protected Admin Routes
	adminGroup := e.Group("/admin")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		SigningKey: auth.SecretKey,
	}
	adminGroup.Use(echojwt.WithConfig(config))

	adminGroup.GET("/dashboard", h.AdminDashboard)
	adminGroup.POST("/restaurants", h.CreateRestaurant)
	adminGroup.DELETE("/restaurants/:id", h.SuspendRestaurant)

	// Restaurant Routes
	restGroup := e.Group("/restaurant")
	restGroup.Use(echojwt.WithConfig(config))
	restGroup.GET("/dashboard", h.RestaurantDashboard)
	restGroup.POST("/menu", h.CreateMenuItem)
	restGroup.PUT("/menu/:id", h.UpdateStock)

	// Guest Routes (Public but logic protected)
	e.POST("/guest/:id/checkin", h.GuestCheckIn)
	e.GET("/guest/:id/menu", h.GuestGetMenu)
	e.POST("/guest/:id/checkout", h.GuestCheckout)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
