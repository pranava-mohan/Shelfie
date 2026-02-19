package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/server/handlers"
)

func InitUser(api fiber.Router) {

	api = api.Group("/user")

	api.Post("/", handlers.GetUser)
	api.Post("/my-books", handlers.GetMyBooks)

}
