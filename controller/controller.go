package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	config "github.com/Naik-Bharat/event-registration/auth"
	"github.com/gofiber/fiber/v2"
)

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

	userData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading user data", err)
	}

	fmt.Println(string(userData))

	err = ctx.SendString("hello")
	return err
}

func GoogleLogin(ctx *fiber.Ctx) error {
	googleConfig := config.Config()
	url := googleConfig.AuthCodeURL("randomstate")

	err := ctx.Redirect(url, fiber.StatusSeeOther)
	return err
}
