package main

import (
	"mobile-order-app/internal/auth"
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/email"
	"mobile-order-app/internal/handlers/api"
	"mobile-order-app/internal/handlers/restaurant"
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

	e.Static("/", "views/restaurant")

	authH := api.NewAuthHandler(database)

	emailSender := &email.LogSender{}
	restH := restaurant.NewHandler(database, emailSender)

	e.POST("/login", authH.Login)

	restGroup := e.Group("/restaurant")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		SigningKey: auth.SecretKey,
	}
	restGroup.Use(echojwt.WithConfig(config))
	// TODO: Role check middleware

	restGroup.GET("/dashboard", restH.Dashboard)
	restGroup.POST("/menu", restH.CreateMenuItem)
	restGroup.PUT("/menu/:id", restH.UpdateStock)
	restGroup.POST("/staff/invite", restH.InviteStaff)
	restGroup.GET("/staff", restH.GetStaff)

	go utils.OpenBrowser("http://localhost:8082")

	e.Logger.Fatal(e.Start(":8082"))
}
