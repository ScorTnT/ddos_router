package controller

import (
	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/gofiber/fiber/v2"
)

func GetIntranetConfig(c *fiber.Ctx) error {
	intranetConfig, err := config.LoadIntranetConfig()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(intranetConfig)
}

func SetIntranetConfig(c *fiber.Ctx) error {
	intranetConfig := config.NewIntranetConfig()

	if err := c.BodyParser(intranetConfig); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := config.ApplyIntranetConfig(intranetConfig); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Config saved!",
		"current_config": intranetConfig,
	})
}
