package controller

import "github.com/gofiber/fiber/v2"

func HookHandler(app *fiber.App) {
	app.Post("/login", Login)
	app.Get("/intranet", GetIntranetConfig)

	app.Get("/network", SearchNetworkConfig)
	app.Get("network/:id", GetIntranetConfig)

	app.Get("/user/:id", GetUser)
	app.Get("/user", SearchUser)
	app.Post("/user", AddUser)
	app.Put("/user", UpdateUser)
	app.Delete("/user", DeleteUser)

	app.Get("/connections", GetConnections)

	app.Get("/", Index)
}
