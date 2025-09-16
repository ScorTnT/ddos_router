package api

import (
	"github.com/ScorTnT/ddos_router/backend/internal"
	"github.com/ScorTnT/ddos_router/backend/internal/protection"
	"github.com/gofiber/fiber/v2"
)

func HandleRoutes(apiGroup fiber.Router, appConfig *internal.Config, protectionManager *protection.Manager) {
	infoGroup := apiGroup.Group("/information")
	infoGroup.Use(ValidateSessionMiddleware)
	infoGroup.Get("/neighbors", GetNeighbors)
	infoGroup.Get("/connections", GetPackets(appConfig))
	infoGroup.Get("/", GetInformation(appConfig))

	configGroup := apiGroup.Group("/config")
	configGroup.Use(ValidateSessionMiddleware)
	configGroup.Get("/wan", GetWANConfig)
	configGroup.Get("/lan", GetLANConfig)
	configGroup.Post("/wan", PostWANConfig)
	configGroup.Post("/lan", PostLANConfig)

	// Blacklist from protection group
	protectionGroup := apiGroup.Group("/protection")
	protectionGroup.Use(ValidateSessionMiddleware)
	protectionGroup.Get("/", GetProtection(protectionManager))
	protectionGroup.Post("/ip/block", PostIPBlock(protectionManager))
	protectionGroup.Post("/ip/unblock", PostIPUnblock(protectionManager))

	// Whitelist from protection group
	whitelistGroup := protectionGroup.Group("/whitelist")
	whitelistGroup.Use(ValidateSessionMiddleware)
	whitelistGroup.Get("/", GetWhitelist(protectionManager))
	whitelistGroup.Post("/add", PostIPAddWhitelist(protectionManager))
	whitelistGroup.Post("/remove", PostIPDeleteWhitelist(protectionManager))

	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", Login)
	authGroup.Post("/logout", Logout)
}
