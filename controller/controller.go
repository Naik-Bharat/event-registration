package controller

import (
	"encoding/json"
	"fmt"
	"io"
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
		println("Error exchanging code", err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		println("Error fetching users details", err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	userData, err := io.ReadAll(res.Body)
	if err != nil {
		println("Error reading user data", err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	user := database.User{}
	err = json.Unmarshal(userData, &user)
	if err != nil {
		println("Error converting user data to struct", err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}
	database.CreateUser(user)

	err = ctx.Redirect("/")
	return err
}

func AddEvent(ctx *fiber.Ctx) error {
	body := new(database.Event)
	err := ctx.BodyParser(body)
	if err != nil {
		println("Cannot parse params", err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}
	err = database.AddEvent(*body)
	if err != nil {
		println(err)
		err = ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}
	err = ctx.SendStatus(fiber.StatusOK)
	return err
}

func BookTicket(ctx *fiber.Ctx) error {
	body := ctx.AllParams()
	fmt.Println(body)
	userID, err := strconv.Atoi(body["user_id"])
	if err != nil {
		println("Error converting user ID to int")
		err = ctx.SendStatus(fiber.StatusBadRequest)
		return err
	}
	eventID, err := strconv.Atoi(body["eventID"])
	if err != nil {
		println("Error converting event ID to int")
		err = ctx.SendStatus(fiber.StatusBadRequest)
		return err
	}
	ticket := database.Ticket{
		UserID:  uint(userID),
		EventID: uint(eventID),
	}
	database.BookTicket(ticket)
	fmt.Println("Booked ticket for user", userID, "event", eventID)

	err = ctx.SendString("hello")
	return err
}

func GoogleLogin(ctx *fiber.Ctx) error {
	googleConfig := config.Config()
	url := googleConfig.AuthCodeURL("randomstate")

	err := ctx.Redirect(url, fiber.StatusSeeOther)
	return err
}
