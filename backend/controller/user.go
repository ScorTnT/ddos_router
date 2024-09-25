package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mchu7797/router_api/database"
	"github.com/mchu7797/router_api/database/models"
)

func SearchUser(c *fiber.Ctx) error {
	username := c.Query("username")

	var user models.User

	result := database.Connection.
		Where("username = ?", username).
		First(&user)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&user)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	result := database.Connection.Find(&user, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&user)
}

func AddUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(user); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	database.Connection.Create(&user)

	return c.Status(201).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	user := new(models.User)
	id := c.Params("id")

	if err := c.BodyParser(user); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	database.Connection.Where("id = ?", id).Updates(&user)

	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := database.Connection.Delete(&user, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}
