package controller

import (
	"github.com/ScorTnT/ddos_router/backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GetArpNowPing(c *fiber.Ctx) error {
	deviceList, err := utils.GetARPNowTable()

	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(200).JSON(&deviceList)
}
