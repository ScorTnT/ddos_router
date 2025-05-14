package controller

import (
	"github.com/ScorTnT/ddos_router/backend/protect_manager"
	"github.com/gofiber/fiber/v2"
)

func HookHandler(app *fiber.App, protectManager *protect_manager.ProtectManager) {
	app.Post("/login", Login)
	app.Get("/logout", Logout)
	app.Post("/update_credentials", UpdateCredentials)

	app.Get("/intranet", GetIntranetConfig)
	app.Post("/intranet", SetIntranetConfig)

	app.Get("/internet", GetInternetConfig)
	app.Post("/internet", SetInternetConfig)

	app.Get("/connections", GetConnections)
	app.Get("/router_info", GetRouterInfo)

	app.Get("/arp", GetArpPing)
	app.Get("/arp_now", GetArpNowPing)

	app.Get("/firewall", GetLockedIPs)
	app.Post("/firewall", LockIP)
	app.Delete("/firewall", UnlockIP)

	protect_api := NewProtectAPI(protectManager)

	app.Get("/protect", protect_api.GetBlockedIPsHandler)
	app.Post("/protect/block", protect_api.BlockIPHandler)
	app.Delete("/protect/unblock", protect_api.UnblockIPHandler)

	app.Get("/", Index)
}
