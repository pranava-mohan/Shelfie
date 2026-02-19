package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/server/handlers"
)

func InitBook(api fiber.Router) {

	api = api.Group("/book")

	api.Post("/create", handlers.CreateBook)
	api.Post("/update", handlers.UpdateBook)
	api.Post("/delete", handlers.DeleteBook)

	api.Post("/get", handlers.GetBook)
	api.Post("/check-in", handlers.CheckInBooks)
	api.Post("/return", handlers.ReturnBooks)

	api.Post("/all", handlers.GetAllBooks)
}
