package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mchu7797/router_api/utils"
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
		return c.Status(503).SendString(err.Error())
	}

	return c.SendString(strconv.FormatBool(utils.CheckAccount(loginData.Username, loginData.Password)))
}
