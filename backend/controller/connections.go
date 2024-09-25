package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mchu7797/router_api/utils"
)

func GetConnections(c *fiber.Ctx) error {
	conntrackEntries, err := utils.GetConntrackEntries()

	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(200).JSON(&conntrackEntries)
}
