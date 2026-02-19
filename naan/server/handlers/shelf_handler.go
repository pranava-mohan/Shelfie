package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"github.com/pranava-mohan/library-automation-pre/naan/server/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateShelfReq struct {
	Address string `json:"address"`
}

func CreateShelf(c *fiber.Ctx) error {
	// check if user is admin
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	data := new(CreateShelfReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shelfCollection := config.GetShelfCollection()

	shelfDetails, shelfErr := shelfCollection.InsertOne(ctx, data)

	if shelfErr != nil {
		return shelfErr
	}

	insertedID := shelfDetails.InsertedID.(bson.ObjectID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "shelf created successfully",
		"shelf_id": insertedID.Hex(),
	})
}

type UpdateShelfReq struct {
	ShelfID string `json:"shelf_id"`
	Address string `json:"address"`
}

func UpdateShelf(c *fiber.Ctx) error {
	// check if user is admin
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	data := new(UpdateShelfReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shelfCollection := config.GetShelfCollection()

	// Convert ShelfID to ObjectID
	shelfID, err := bson.ObjectIDFromHex(data.ShelfID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid shelf ID",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"address": data.Address,
		},
	}

	_, err = shelfCollection.UpdateOne(ctx, bson.M{"_id": shelfID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update shelf",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shelf updated successfully",
	})
}

type DeleteShelfReq struct {
	ShelfID string `json:"shelf_id"`
}

func DeleteShelf(c *fiber.Ctx) error {
	// check if user is admin
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	data := new(DeleteShelfReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shelfCollection := config.GetShelfCollection()
	// Convert ShelfID to ObjectID
	shelfID, err := bson.ObjectIDFromHex(data.ShelfID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid shelf ID",
		})
	}
	_, err = shelfCollection.DeleteOne(ctx, bson.M{"_id": shelfID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete shelf",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shelf deleted successfully",
	})
}

func GetAllShelves(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shelfCollection := config.GetShelfCollection()

	cursor, err := shelfCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch shelves",
		})
	}
	defer cursor.Close(ctx)

	var shelves []models.Shelf
	if err = cursor.All(ctx, &shelves); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode shelves",
		})
	}

	return c.Status(fiber.StatusOK).JSON(shelves)
}
