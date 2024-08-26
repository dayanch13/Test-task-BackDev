package main

import (
	"github.com/gofiber/fiber/v2"
	"jwt-auth/controllers"
	"jwt-auth/db"
)

func main() {
	db.DBconn()
	router := fiber.New()
	controllers.Setup(router)
	err := router.Listen(":8080")
	if err != nil {
		return
	}
}
