package controller

import (
	"github.com/ScorTnT/ddos_router/backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GetRouterInfo(c *fiber.Ctx) error {
	info, err := utils.GetRouterInfo()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(info)
}
