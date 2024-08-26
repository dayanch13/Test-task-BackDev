package controllers

import (
	"github.com/gofiber/fiber/v2"
	"jwt-auth/api"
)

func Setup(app *fiber.App) {
	token := app.Group("/api")
	token.Post("/login", api.Login)
	token.Post("/renew", api.Renew)
}
