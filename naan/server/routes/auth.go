package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/server/handlers"
)

func InitAuth(api fiber.Router) {
	// api.Get("/login/dauth", handlers.LoginDauth)
	// api.Get("/auth/dauth", handlers.DauthCallback)
	api.Get("/login/google", handlers.LoginGoogle)
	api.Get("/auth/google", handlers.GoogleCallback)

	api.Post("/login/admin", handlers.LoginAdmin)

}
