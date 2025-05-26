package api

import (
	"github.com/ScorTnT/ddos_router/backend/internal/protection"
	"github.com/gofiber/fiber/v2"
)

func GetProtection(protectionManager *protection.ProtectionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blockedIP := protectionManager.SnapshotBlocks()

		if len(blockedIP) == 0 {
			return RespondWithError(c, fiber.StatusNotFound, "Blocked IP is Empty")
		}

		return RespondWithJSON(c, fiber.StatusOK, blockedIP)
	}
}

func PostIPBlock(protectionManager *protection.ProtectionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.Query("ip")
		if ip == "" {
			return RespondWithError(c, fiber.StatusBadRequest, "IP address is required")
		}

		protectionManager.Add(ip)

		return RespondWithCustomMessage(c, fiber.StatusOK, "IP address added from firewall rules")
	}
}

func PostIPUnblock(protectionManager *protection.ProtectionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.Query("ip")
		if ip == "" {
			return RespondWithError(c, fiber.StatusBadRequest, "IP address is required")
		}

		protectionManager.Delete(ip)

		return RespondWithCustomMessage(c, fiber.StatusOK, "IP address deleted from firewall rules")
	}
}
