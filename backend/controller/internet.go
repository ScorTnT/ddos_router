package controller

import (
	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/gofiber/fiber/v2"
)

func GetInternetConfig(c *fiber.Ctx) error {
	internetConfig, err := config.LoadInternetConfig()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(internetConfig)
}

func SetInternetConfig(c *fiber.Ctx) error {
	internetConfig := config.InternetConfig{}

	if err := c.BodyParser(internetConfig); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := config.ApplyInternetConfig(&internetConfig); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Config saved!",
		"current_config": internetConfig,
	})
}
