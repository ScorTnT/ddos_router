package api

import (
	"github.com/ScorTnT/ddos_router/backend/internal"
	"github.com/ScorTnT/ddos_router/backend/internal/router"
	"github.com/gofiber/fiber/v2"
)

func GetNeighbors(c *fiber.Ctx) error {
	neighbors, err := router.GetOnlineNeighbors()
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	return RespondWithJSON(c, fiber.StatusOK, neighbors)
}

func GetPackets(appConfig *internal.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		packets, err := router.GetPackets(appConfig)
		if err != nil {
			return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
		}

		return RespondWithJSON(c, fiber.StatusOK, packets)
	}
}

func GetInformation(appConfig *internal.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		information, err := router.GetInformation(appConfig)
		if err != nil {
			return RespondWithError(c, fiber.StatusInternalServerError, err.Error())
		}

		return RespondWithJSON(c, fiber.StatusOK, information)
	}
}
