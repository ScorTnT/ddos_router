package controller

import "github.com/gofiber/fiber/v2"

func GetIntranetConfig(c *fiber.Ctx) error {
	return c.SendString("Intranet Handler.")
}

func AddIntranetConfig(c *fiber.Ctx) error {
	return nil
}

func UpdateIntranetConfig(c *fiber.Ctx) error {
	return nil
}

func DeleteIntranetConfig(c *fiber.Ctx) error {
	return nil
}
