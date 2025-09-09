package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ScorTnT/ddos_router/backend/internal"
	"github.com/ScorTnT/ddos_router/backend/internal/api"
	"github.com/ScorTnT/ddos_router/backend/internal/firewall"
	"github.com/ScorTnT/ddos_router/backend/internal/protection"
	"github.com/ScorTnT/ddos_router/backend/internal/snort"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	appConfig, err := internal.LoadConfig()
	if err != nil {
		log.Printf("[WARNING] Failed to load config: %v", err)
		log.Println("[INFO] Using default configuration")
		appConfig = &internal.Config{
			WebAPI: internal.WebAPIConfig{
				ListenAddr: ":2024",
			},
			Interface: internal.InferfaceConfig{
				WANInterfaceName: "eth0",
				LANInterfaceName: "eth1",
			},
			Snort: internal.SnortAlertScannerConfig{
				SnortLogPath: "/var/log/snort/alert",
				BufferSize:   1024,
			},
			Protection: internal.ProtectionConfig{
				ProtectionTTL: 30,
				RefreshTick:   1000,
			},
		}
	}

	log.Printf("[INFO] Loaded configuration. API listen address: %s", appConfig.WebAPI.ListenAddr)

	if err := firewall.InitFirewall(); err != nil {
		log.Fatalf("[ERROR] Failed to initialize firewall: %v", err)
	}

	defer func() {
		if err := firewall.CleanupFirewall(); err != nil {
			log.Printf("[WARNING] Failed to clean up firewall: %v", err)
		}
	}()

	alertScanner := snort.NewSnortScanner(
		appConfig.Snort.SnortLogPath,
		appConfig.Snort.BufferSize,
		nil,
	)

	protectionManager, err := protection.NewProtectManager(
		time.Duration(appConfig.Protection.ProtectionTTL)*time.Second,
		time.Duration(appConfig.Protection.RefreshTick)*time.Millisecond,
		alertScanner,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create protection manager: %v", err)
	}

	protectionManager.Start()
	log.Println("[INFO] Protection manager started")
	defer func() {
		protectionManager.Stop()
		log.Println("[INFO] Protection manager stopped")
	}()

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				message := "Internal Server Error"

				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
					message = e.Message
				}

				log.Printf("[ERROR] %s: %s", message, err)
				return api.RespondWithError(c, code, message)
			},
		},
	)

	app.Use(
		cors.New(
			cors.Config{
				AllowOrigins: "*",
				AllowHeaders: "Origin, Content-Type, Accept, " + api.SessionIDHeader,
				AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
			},
		),
	)
	app.Use(api.RequestLogger())

	api.HandleRoutes(app, appConfig, protectionManager)

	app.Get(
		"/",
		func(c *fiber.Ctx) error {
			return api.RespondWithJSON(
				c,
				fiber.StatusOK,
				fiber.Map{
					"message": "DDOS Router API",
				},
			)
		},
	)

	go func() {
		listenAddr := appConfig.WebAPI.ListenAddr
		log.Printf("[INFO] Starting web API server on %s", listenAddr)
		if err := app.Listen(listenAddr); err != nil {
			log.Fatalf("[ERROR] Failed to start web API server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	receivedSignal := <-quit
	log.Printf("[INFO] Received signal: %s", receivedSignal)
	log.Println("[INFO] Shutting down web API server...")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownContext); err != nil {
		log.Fatalf("Server Shutdown(Force): %v", err)
	}

	log.Println("[INFO] Web API server shut down gracefully")
}
