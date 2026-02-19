package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/server/handlers"
)

func InitShelf(api fiber.Router) {
	api = api.Group("/shelf")

	api.Post("/create", handlers.CreateShelf)
	api.Post("/update", handlers.UpdateShelf)
	api.Post("/delete", handlers.DeleteShelf)
	api.Post("/all", handlers.GetAllShelves)
}
