package main

import (
	"mobile-order-app/internal/db"
	"mobile-order-app/internal/handlers/guest"
	"mobile-order-app/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database := db.Init()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "views/guest")

	guestH := guest.NewHandler(database)

	e.POST("/guest/:id/checkin", guestH.CheckIn)
	e.GET("/guest/:id/menu", guestH.GetMenu)
	e.POST("/guest/:id/checkout", guestH.Checkout)

	go utils.OpenBrowser("http://localhost:8080")

	e.Logger.Fatal(e.Start(":8080"))
}
