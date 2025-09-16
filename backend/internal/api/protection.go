package api

import (
	"net"
	"strconv"
	"strings"

	"github.com/ScorTnT/ddos_router/backend/internal/protection"
	"github.com/gofiber/fiber/v2"
)

func GetProtection(protectionManager *protection.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blockedIP := protectionManager.SnapshotBlocks()

		if len(blockedIP) == 0 {
			return RespondWithError(c, fiber.StatusNotFound, "Blocked IP is Empty")
		}

		return RespondWithJSON(c, fiber.StatusOK, blockedIP)
	}
}

func PostIPBlock(protectionManager *protection.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := strings.TrimSpace(c.Query("ip"))
		if ip == "" {
			return RespondWithError(c, fiber.StatusBadRequest, "IP address is required")
		}

		parsed := net.ParseIP(ip)
		if parsed == nil {
			return RespondWithError(c, fiber.StatusBadRequest, "Invalid IP address")
		}

		isPermanent, err := strconv.ParseBool(c.Query("is_permanent", "false"))

		if err != nil {
			return RespondWithError(c, fiber.StatusBadRequest, "Invalid value for 'is_permanent'")
		}

		protectionManager.Add(parsed.String(), isPermanent)

		return RespondWithCustomMessage(c, fiber.StatusOK, "IP address added from firewall rules")
	}
}

func PostIPUnblock(protectionManager *protection.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := strings.TrimSpace(c.Query("ip"))
		if ip == "" {
			return RespondWithError(c, fiber.StatusBadRequest, "IP address is required")
		}

		parsed := net.ParseIP(ip)
		if parsed == nil {
			return RespondWithError(c, fiber.StatusBadRequest, "Invalid IP address")
		}

		protectionManager.Delete(parsed.String())

		return RespondWithCustomMessage(c, fiber.StatusOK, "IP address deleted from firewall rules")
	}
}
