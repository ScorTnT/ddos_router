package controller

import (
	"github.com/ScorTnT/ddos_router/backend/protect_manager"
	"github.com/gofiber/fiber/v2"
)

type ProtectAPI struct {
	ProtectManager *protect_manager.ProtectManager
}

func NewProtectAPI(protectManager *protect_manager.ProtectManager) *ProtectAPI {
	return &ProtectAPI{
		ProtectManager: protectManager,
	}
}

func (api *ProtectAPI) BlockIPHandler(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IP address is required"})
	}

	api.ProtectManager.Add(ip)

	return c.JSON(fiber.Map{"message": "IP block request sended successfully"})
}

func (api *ProtectAPI) UnblockIPHandler(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IP address is required"})
	}

	api.ProtectManager.Delete(ip)

	return c.JSON(fiber.Map{"message": "IP unblock request sended successfully"})
}

func (api *ProtectAPI) GetBlockedIPsHandler(c *fiber.Ctx) error {
	blockedIPs := api.ProtectManager.SnapshotBlocks()
	return c.JSON(blockedIPs)
}
