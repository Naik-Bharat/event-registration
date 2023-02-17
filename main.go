package main

import (
	"github.com/Naik-Bharat/event-registration/controller"
	"github.com/Naik-Bharat/event-registration/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Database
	database.ConnectDB()
	db := database.DB
	database.AutoMigrate(db)

	// web app
	app := fiber.New()
	app.Get("/auth/google/login", controller.GoogleLogin)
	app.Get("/auth/google/callback", controller.GoogleCallback)
	app.Get("/", controller.Index)

	app.Post("/api/add_event", controller.AddEvent)
	app.Listen(":8080")
}
