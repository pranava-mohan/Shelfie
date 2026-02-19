package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/server/handlers"
)

func InitKiosk(api fiber.Router) {
	api = api.Group("/kiosk")

	api.Post("/create", handlers.CreateKiosk)
	api.Post("/delete", handlers.DeleteKiosk)
	api.Get("/list", handlers.ListKiosks)
	api.Get("/:kiosk_name", handlers.KioskAuth)
}
