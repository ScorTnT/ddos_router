package controller

import (
	"strconv"

	"github.com/ScorTnT/ddos_router/backend/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *fiber.Ctx) error {
	loginData := new(LoginRequest)

	if err := c.BodyParser(loginData); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendString(strconv.FormatBool(utils.CheckAccount(loginData.Username, loginData.Password)))
}
