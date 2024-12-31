package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/ScorTnT/ddos_router/backend/controller"
	"github.com/ScorTnT/ddos_router/backend/database"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	controller.HookHandler(app)

	var err error

	routerConfig, err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	err = database.Init("secure_router.db", routerConfig)

	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":2024"))
}
