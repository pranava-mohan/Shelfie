package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"github.com/pranava-mohan/library-automation-pre/naan/server/services"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GetUserReq struct {
	ID string `json:"id"`
}

func GetUser(c *fiber.Ctx) error {

	data := new(GetUserReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	fmt.Println(data.ID)

	userObjID, err := bson.ObjectIDFromHex(data.ID)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	userCollection := config.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.PublicUser
	err = userCollection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func GetProfile(c *fiber.Ctx) error {
	userIDString := services.GetUserID(c)
	if userIDString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userObjID, err := bson.ObjectIDFromHex(userIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userCollection := config.GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

/*
output:

list of books with their id, title, author, publisher, isbn, and the date user took the book (issued_at), don't return books if he returned it
*/

type BorrowedBookResult struct {
	Book     models.PublicBook `bson:"book_details" json:"book_details"`
	IssuedAt time.Time         `bson:"issued_at" json:"issued_at"`
}

func GetMyBooks(c *fiber.Ctx) error {

	if !services.IsNormalUser(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	userIDString := services.GetUserID(c)
	userID, err := bson.ObjectIDFromHex(userIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	historyCollection := config.GetHistoryCollection()

	pipeline := mongo.Pipeline{
		// Stage 1: Filter History to find active issues for this user
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "returned_at", Value: nil}, // Matches if field is null or does not exist
		}}},

		// Stage 2: Join with the 'books' collection
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},         // The name of your book collection
			{Key: "localField", Value: "book_id"}, // Field in History
			{Key: "foreignField", Value: "_id"},   // Field in Book
			{Key: "as", Value: "book_details"},    // Output array field
		}}},

		// Stage 3: Unwind the 'book_details' array
		{{Key: "$unwind", Value: "$book_details"}},

		// Stage 4: Project to ensure clean output
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},          // Exclude the history ID
			{Key: "issued_at", Value: 1},    // Keep issued_at from history
			{Key: "book_details", Value: 1}, // Keep the joined book object
		}}},
	}

	// Execute the pipeline
	cursor, err := historyCollection.Aggregate(c.Context(), pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(c.Context())

	// Decode results
	var results []BorrowedBookResult
	if err := cursor.All(c.Context(), &results); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(results)
}
