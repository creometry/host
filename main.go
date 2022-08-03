package main

import (
	"api/controller"

	"github.com/gofiber/fiber/v2"
)

func Routers(app *fiber.App) {
	app.Get("/payment", controller.Payment)
}
func main() {
	app := fiber.New()
	Routers(app)
	app.Listen(":3000")
}
