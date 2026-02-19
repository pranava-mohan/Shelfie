package services

import "github.com/gofiber/fiber/v2"

func IsAdmin(c *fiber.Ctx) bool {
	userStatus := c.Locals("user_type")
	return userStatus == "admin"
}

func IsKiosk(c *fiber.Ctx) bool {
	userStatus := c.Locals("user_type")
	return userStatus == "kiosk"
}

func IsNormalUser(c *fiber.Ctx) bool {
	userStatus := c.Locals("user_type")
	return userStatus == "normal"
}

func GetUserID(c *fiber.Ctx) string {
	return c.Locals("user_id").(string)
}
