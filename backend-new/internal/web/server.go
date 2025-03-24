package web

import (
	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return app
}
