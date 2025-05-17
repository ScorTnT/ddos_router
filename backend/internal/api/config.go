package api

import (
	"github.com/ScorTnT/ddos_router/backend/internal/router/config"
	"github.com/gofiber/fiber/v2"
)

func GetWANConfig(c *fiber.Ctx) error {
	wanConfig, err := config.LoadWANConfig()
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return RespondWithJSON(c, fiber.StatusOK, wanConfig)
}

func GetLANConfig(c *fiber.Ctx) error {
	lanConfig, err := config.LoadLANConfig()
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return RespondWithJSON(c, fiber.StatusOK, lanConfig)
}

func PostWANConfig(c *fiber.Ctx) error {
	wanConfig, err := config.LoadWANConfig()
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := c.BodyParser(&wanConfig); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := config.ApplyWANConfig(wanConfig); err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return RespondWithJSON(c, fiber.StatusOK, wanConfig)
}

func PostLANConfig(c *fiber.Ctx) error {
	lanConfig, err := config.LoadLANConfig()
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := c.BodyParser(&lanConfig); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := config.ApplyLANConfig(lanConfig); err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return RespondWithJSON(c, fiber.StatusOK, lanConfig)
}
