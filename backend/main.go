package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	"github.com/ScorTnT/ddos_router/backend/controller"
	"github.com/ScorTnT/ddos_router/backend/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	controller.HookHandler(app)

	err := database.Connect("example.db")

	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":2024"))
}
