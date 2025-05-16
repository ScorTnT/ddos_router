package api

import (
	"github.com/gofiber/fiber/v2"
)

func HandleRoutes(apiGroup fiber.Router) {
	infoGroup := apiGroup.Group("/information")
	infoGroup.Use(ValidateSessionMiddleware)
	infoGroup.Get("/neighbors", GetNeighbors)
	infoGroup.Get("/connections", GetConnections)
	infoGroup.Get("/", GetInformation)

	configGroup := apiGroup.Group("/config")
	configGroup.Use(ValidateSessionMiddleware)
	configGroup.Get("/wan", GetWANConfig)
	configGroup.Get("/lan", GetLANConfig)
	configGroup.Post("/wan", PostWANConfig)
	configGroup.Post("/lan", PostLANConfig)

	protectionGroup := apiGroup.Group("/protection")
	protectionGroup.Use(ValidateSessionMiddleware)
	protectionGroup.Get("/", GetProtection)
	protectionGroup.Post("/ip/block", PostIPBlock)
	protectionGroup.Post("/ip/unblock", PostIPUnblock)

	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", Login)
	authGroup.Post("/logout", Logout)
}
