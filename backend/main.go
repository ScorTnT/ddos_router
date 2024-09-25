package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mchu7797/router_api/controller"
	"github.com/mchu7797/router_api/database"
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
