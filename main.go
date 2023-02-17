package main

import (
	"github.com/Naik-Bharat/event-registration/controller"
	"github.com/Naik-Bharat/event-registration/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db := database.ConnectDB()
	database.AutoMigrate(db)
	app := fiber.New()
	app.Get("/auth/google/login", controller.GoogleLogin)
	app.Get("/auth/google/callback", controller.GoogleCallback)
	app.Listen(":8080")
}
