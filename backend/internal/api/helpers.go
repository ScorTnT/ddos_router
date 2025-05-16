package api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type StandardResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()
		log.Printf("[WebAPI] [%s] %s %s - %d %s", c.IP(), c.Method(), c.Path(), c.Response().StatusCode(), stop.Sub(start))
		return err
	}
}

func RespondWithError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(StandardResponse{
		Status:  "error",
		Message: message,
	})
}

func RespondWithJSON(c *fiber.Ctx, code int, data interface{}) error {
	return c.Status(code).JSON(StandardResponse{
		Status: "success",
		Data:   data,
	})
}
