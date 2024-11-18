package controller

import "github.com/gofiber/fiber/v2"

func HookHandler(app *fiber.App) {
	app.Post("/login", Login)
	app.Get("/logout", Logout)
	app.Post("/update_credentials", UpdateCredentials)

	app.Get("/intranet", GetIntranetConfig)
	app.Post("/intranet", SetIntranetConfig)

	app.Get("/internet", GetInternetConfig)
	app.Post("/internet", SetInternetConfig)

	app.Get("/connections", GetConnections)
	app.Get("/router_info", GetRouterInfo)

	app.Get("/", Index)
}
