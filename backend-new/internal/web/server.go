package web

import (
	"openwrt-ips/internal/web/controller"

	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	app := fiber.New()
	controller.HookHandler(app)

	return app
}
