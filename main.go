package main

import (
	"github.com/Naik-Bharat/event-registration/controller"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/auth/google/login", controller.GoogleLogin)
	app.Get("/auth/google/callback", controller.GoogleCallback)
	app.Listen(":8080")
}
