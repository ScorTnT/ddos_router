package controller

import "github.com/gofiber/fiber/v2"

func HookHandler(app *fiber.App) {
	app.Post("/login", Login)

	app.Get("/intranet", GetIntranetConfig)
	app.Post("/intranet", SetIntranetConfig)

	app.Get("/internet", GetInternetConfig)
	app.Post("/internet", SetInternetConfig)

	app.Get("/user/:id", GetUser)
	app.Get("/user", SearchUser)
	app.Post("/user", AddUser)
	app.Put("/user", UpdateUser)
	app.Delete("/user", DeleteUser)

	app.Get("/connections", GetConnections)

	app.Get("/", Index)
}
