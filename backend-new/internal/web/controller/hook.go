package controller

import "github.com/gofiber/fiber/v2"

func HookHandler(app *fiber.App) {
	app.Get("/arp", GetArpData)
}
