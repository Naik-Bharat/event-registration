package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	config "github.com/Naik-Bharat/event-registration/auth"
	"github.com/Naik-Bharat/event-registration/database"
	"github.com/gofiber/fiber/v2"
)

func Index(ctx *fiber.Ctx) error {
	err := ctx.SendString("hello")
	return err
}

func GoogleCallback(ctx *fiber.Ctx) error {
	googleConfig := config.Config()
	token, err := googleConfig.Exchange(ctx.Context(), ctx.Query("code"))
	if err != nil {
		log.Fatal("Error exchanging code", err)
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Fatal("Error fetching users details", err)
	}

	userData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading user data", err)
	}

	user := database.User{}
	err = json.Unmarshal(userData, &user)
	if err != nil {
		log.Fatal("Error converting user data to struct", err)
	}
	database.CreateUser(user)

	err = ctx.Redirect("/")
	return err
}

func AddEvent(ctx *fiber.Ctx) error {
	body := new(database.Event)
	err := ctx.BodyParser(body)
	if err != nil {
		log.Fatal("Cannot parse params", err)
	}
	database.AddEvent(*body)
	err = ctx.SendStatus(fiber.StatusOK)
	return err
}

func BookTicket(ctx *fiber.Ctx) error {
	body := ctx.AllParams()
	fmt.Println(body)
	userID, err := strconv.Atoi(body["user_id"])
	eventID, err := strconv.Atoi(body["eventID"])
	ticket := database.Ticket{
		UserID:  uint(userID),
		EventID: uint(eventID),
	}
	database.BookTicket(ticket)
	err = ctx.SendString("hello")
	return err
}

func GoogleLogin(ctx *fiber.Ctx) error {
	googleConfig := config.Config()
	url := googleConfig.AuthCodeURL("randomstate")

	err := ctx.Redirect(url, fiber.StatusSeeOther)
	return err
}
