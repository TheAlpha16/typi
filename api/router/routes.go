package router

import (
	"github.com/TheAlpha16/typi/api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/ping", handler.Ping)
	api.Get("/videos", handler.GetVideos)
}
