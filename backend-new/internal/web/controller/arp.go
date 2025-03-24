package controller

import (
	"openwrt-ips/internal/analyzer"
	"openwrt-ips/internal/config"

	"github.com/gofiber/fiber/v2"
)

func GetArpData(c *fiber.Ctx) error {
	cfg := config.GetConfig()
	device_list, err := analyzer.GetARPTable(cfg)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get arp data!")
	}

	return c.JSON(device_list)
}
