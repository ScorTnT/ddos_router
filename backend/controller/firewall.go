package controller

import (
	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/gofiber/fiber/v2"
)

// IPInfo는 방화벽에서 사용하는 IP 정보를 담는 구조체입니다
type IPInfo struct {
	IPAddress string `json:"ip_address"`
}

// GetLockedIPs는 현재 차단된 IP 목록을 조회합니다
func GetLockedIPs(c *fiber.Ctx) error {
	blockedIPs, err := config.GetBlockedIPs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 문자열 배열을 IPInfo 구조체 배열로 변환
	ipInfoList := make([]IPInfo, len(blockedIPs))
	for i, ip := range blockedIPs {
		ipInfoList[i] = IPInfo{IPAddress: ip}
	}

	return c.JSON(fiber.Map{
		"blocked_ips": ipInfoList,
		"count":       len(ipInfoList),
	})
}

// LockIP는 특정 IP를 차단 목록에 추가합니다
func LockIP(c *fiber.Ctx) error {
	var ipInfo IPInfo

	if err := c.BodyParser(&ipInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if ipInfo.IPAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "IP address is empty",
		})
	}

	if err := config.BlockIP(ipInfo.IPAddress); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "IP successfully blocked",
		"ip_info": ipInfo,
	})
}

// UnlockIP는 차단된 IP를 차단 목록에서 제거합니다
func UnlockIP(c *fiber.Ctx) error {
	var ipInfo IPInfo

	if err := c.BodyParser(&ipInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if ipInfo.IPAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "IP address is empty",
		})
	}

	if err := config.UnblockIP(ipInfo.IPAddress); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "IP successfully unblocked",
		"ip_info": ipInfo,
	})
}
