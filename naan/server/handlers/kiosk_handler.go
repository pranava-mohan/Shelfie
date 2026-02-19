package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"github.com/pranava-mohan/library-automation-pre/naan/server/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type KioskReq struct {
	Name string `json:"name"`
}

func CreateKiosk(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}
	var data KioskReq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	kioskCollection := config.GetKioskCollection()
	newKiosk := models.Kiosk{
		Name: data.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, insertErr := kioskCollection.InsertOne(ctx, newKiosk)

	if insertErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create kiosk",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "kiosk created successfully",
		"kiosk_name": newKiosk.Name,
	})

}

func DeleteKiosk(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}
	var data KioskReq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	kioskCollection := config.GetKioskCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, deleteErr := kioskCollection.DeleteOne(ctx, models.Kiosk{
		Name: data.Name,
	})
	if deleteErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete kiosk",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "kiosk deleted successfully",
		"kiosk_name": data.Name,
	})
}

func ListKiosks(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized kiosk",
		})
	}
	kioskCollection := config.GetKioskCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := kioskCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch kiosks",
		})
	}
	var kiosks []models.Kiosk
	if err := cursor.All(ctx, &kiosks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse kiosks",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"kiosks": kiosks,
	})
}

func KioskAuth(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized kiosk",
		})
	}
	kioskName := c.Params("kiosk_name")
	kioskCollection := config.GetKioskCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var kiosk models.Kiosk
	err := kioskCollection.FindOne(ctx, bson.M{"name": kioskName}).Decode(&kiosk)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized kiosk",
		})
	}
	claims := jwt.MapClaims{
		"id":   kiosk.Name,
		"type": "kiosk",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "kiosk authorized",
		"token":   t,
	})
}
