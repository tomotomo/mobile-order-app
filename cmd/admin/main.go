package main

import (
	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/email"
	"mobile-order-app/internal/handlers/admin"
	"mobile-order-app/internal/handlers/api"
	"mobile-order-app/internal/utils"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database := db.Init()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static Files for Admin UI
	e.Static("/", "views/admin")

	// Dependencies
	authH := api.NewAuthHandler(database)
	
	// Use LogSender (Trap) by default
	emailSender := &email.LogSender{}
	// For Mailtrap, uncomment and fill:
	// emailSender := &email.SmtpSender{Host: "sandbox.smtp.mailtrap.io", Port: "2525", Username: "...", Password: "..."}

	adminH := admin.NewHandler(database, emailSender)

	// Public Routes
	e.POST("/login", authH.Login)

	// Protected Routes
	adminGroup := e.Group("/admin")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		SigningKey: auth.SecretKey,
	}
	adminGroup.Use(echojwt.WithConfig(config))
	// TODO: Middleware to check strict Role == Admin

	adminGroup.GET("/dashboard", adminH.Dashboard)
	adminGroup.POST("/restaurants", adminH.CreateRestaurant)
	adminGroup.DELETE("/restaurants/:id", adminH.SuspendRestaurant)

	go utils.OpenBrowser("http://localhost:8081")

	e.Logger.Fatal(e.Start(":8081"))
}
