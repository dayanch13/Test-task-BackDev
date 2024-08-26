package api

import (
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"jwt-auth/db"
	"jwt-auth/models"
	"jwt-auth/services"
	"time"
)

func Login(c *fiber.Ctx) error {
	ip := c.IP()
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.Users
	var session models.Sessions

	db.DB.Where("email = ?", data["email"]).First(&user) //Check the email is present in the DB

	if user.ID == 0 { //If the ID return is '0' then there is no such email present in the DB
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	accessToken, accessClaims, err := services.CreateToken(user.ID, user.Email, ip, 15*time.Minute)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating access token",
		})
	}
	refreshToken, refreshClaims, err := services.CreateToken(user.ID, user.Email, ip, 15*time.Minute)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating refresh token",
		})
	}
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error hashing token",
		})
	}
	session.ID = refreshClaims.RegisteredClaims.ID
	session.UserEmail = refreshClaims.UserEmail
	session.UserIP = refreshClaims.UserIP
	session.RefreshToken = string(hashedToken)

	db.DB.Create(&session)

	return c.JSON(fiber.Map{
		"message":                 "success",
		"access_token":            accessToken,
		"access_token_expire_at":  accessClaims.RegisteredClaims.ExpiresAt,
		"refresh_token":           refreshToken,
		"refresh_token_expire_at": refreshClaims.RegisteredClaims.ExpiresAt,
	})
}

func Renew(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	refreshClaim, err := services.VerifyToken(data["refresh_token"])
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error verifying refresh claims",
		})
	}
	var session models.Sessions

	db.DB.Where("id = ?", refreshClaim.RegisteredClaims.ID).First(&session)
	if session.ID == "" { //If the ID return is '0' then there is no such email present in the DB
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "session not found",
		})
	}
	if session.Attemp > 3 {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating access token you not have attemps to create access token",
		})
	}
	if session.UserEmail != refreshClaim.UserEmail {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating access token you email don't match to create access token",
		})
	}
	if session.UserIP != refreshClaim.UserIP {
		err := services.SendEmail(session.UserEmail, "IP address is changed", "Your IP address has changed")
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "error sending email",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating access token you ip don't match to create access token",
		})

	}
	accessToken, accessClaims, err := services.CreateToken(refreshClaim.ID, refreshClaim.UserEmail, refreshClaim.UserIP, 15*time.Minute)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error creating access token",
		})
	}
	return c.JSON(fiber.Map{
		"message":                "success",
		"access_token":           accessToken,
		"access_token_expire_at": accessClaims.RegisteredClaims.ExpiresAt,
	})
}
