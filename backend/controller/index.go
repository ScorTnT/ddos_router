package controller

import "github.com/gofiber/fiber/v2"

func Index(c *fiber.Ctx) error {
	return c.SendString("Router API")
}
