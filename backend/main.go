package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/ScorTnT/ddos_router/backend/controller"
	"github.com/ScorTnT/ddos_router/backend/database"
	"github.com/ScorTnT/ddos_router/backend/protect_manager"
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

	if err = config.InitFirewall(); err != nil {
		panic(err)
	}

	defer func() {
		if err := config.CleanupFirewall(); err != nil {
			fmt.Println(err)
		}
	}()

	snortScanner := protect_manager.NewSnortScanner(
		"/etc/snort/alert_json.txt",
		128,
		nil,
	)

	protectManager, err := protect_manager.NewProtectManager(
		30*time.Second,
		5*time.Second,
		snortScanner,
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	protectManager.Start()
	defer protectManager.Stop()

	log.Fatal(app.Listen(":2028"))
}
