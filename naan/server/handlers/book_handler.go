package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"github.com/pranava-mohan/library-automation-pre/naan/server/services"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CreateBookReq struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
	Genre     string `json:"genre"`
	ShelfID   string `json:"shelf_id"`
	Row       int    `json:"row"`
	Column    int    `json:"column"`
}

func CreateBook(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	bookCollection := config.GetBookCollection()
	data := new(CreateBookReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	shelfIdObjID, shelfIDErr := bson.ObjectIDFromHex(data.ShelfID)
	if shelfIDErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid shelf ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingBook bson.M
	bookCollection.FindOne(ctx, bson.M{
		"shelf_id": shelfIdObjID,
		"row":      data.Row,
		"column":   data.Column,
	}).Decode(&existingBook)
	if existingBook != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "a book already exists at the specified row and column in the shelf",
		})
	}

	newBook := bson.M{
		"title":     data.Title,
		"author":    data.Author,
		"publisher": data.Publisher,
		"isbn":      data.ISBN,
		"genre":     data.Genre,
		"shelf_id":  shelfIdObjID,
		"added_at":  time.Now(),
		"row":       data.Row,
		"column":    data.Column,
	}
	result, insertErr := bookCollection.InsertOne(ctx, newBook)
	if insertErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book created successfully",
		"book_id": result.InsertedID.(bson.ObjectID).Hex(),
	})
}

type UpdateBookReq struct {
	BookID    string `json:"book_id"`
	Title     string `json:"title,omitempty"`
	Author    string `json:"author,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	ISBN      string `json:"isbn,omitempty"`
	Genre     string `json:"genre,omitempty"`
	ShelfID   string `json:"shelf_id,omitempty"`
	Row       int    `json:"row,omitempty"`
	Column    int    `json:"column,omitempty"`
}

func UpdateBook(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	data := new(UpdateBookReq)
	bookCollection := config.GetBookCollection()
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	bookID, bookIDErr := bson.ObjectIDFromHex(data.BookID)
	if bookIDErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid book ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingBook models.Book
	bookCollection.FindOne(ctx, bson.M{
		"book_id": bookID,
	}).Decode(&existingBook)
	if existingBook.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "book not found",
		})
	}

	update := bson.M{
		"row":    data.Row,
		"column": data.Column,
	}
	if data.Title != "" {
		update["title"] = data.Title
	}
	if data.Author != "" {
		update["author"] = data.Author
	}
	if data.Publisher != "" {
		update["publisher"] = data.Publisher
	}
	if data.ISBN != "" {
		update["isbn"] = data.ISBN
	}
	if data.Genre != "" {
		update["genre"] = data.Genre
	}
	if data.ShelfID != "" {
		shelfIdObjID, shelfIDErr := bson.ObjectIDFromHex(data.ShelfID)
		if shelfIDErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid shelf ID",
			})
		}
		update["shelf_id"] = shelfIdObjID
	}
	if data.Row != 0 {
		update["row"] = data.Row
	}
	if data.Column != 0 {
		update["column"] = data.Column
	}
	// check if smthg is already there in that specified row/col in shelf
	if shelfID, ok := update["shelf_id"]; ok {
		var bookAtPosition bson.M
		bookCollection.FindOne(c.Context(), bson.M{
			"shelf_id": shelfID,
			"row":      update["row"],
			"column":   update["column"],
			"_id":      bson.M{"$ne": bookID},
		}).Decode(&bookAtPosition)
		if bookAtPosition != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "another book already exists at the specified row and column in the shelf",
			})
		}
	}

	_, updateErr := bookCollection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{
		"$set": update,
	})
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book updated successfully",
	})
}

type DeleteBookReq struct {
	BookID string `json:"book_id"`
}

func DeleteBook(c *fiber.Ctx) error {
	if !services.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	data := new(DeleteBookReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bookCollection := config.GetBookCollection()
	bookID, bookIDErr := bson.ObjectIDFromHex(data.BookID)
	if bookIDErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid book ID",
		})
	}
	_, deleteErr := bookCollection.DeleteOne(ctx, bson.M{"_id": bookID})
	if deleteErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book deleted successfully",
	})
}

type GetBookReq struct {
	BookID string `json:"book_id"`
}

func GetBook(c *fiber.Ctx) error {

	if !services.IsKiosk(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	data := new(GetBookReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	bookID, bookIDErr := bson.ObjectIDFromHex(data.BookID)
	if bookIDErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid book ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bookCollection := config.GetBookCollection()
	var book models.PublicBook
	err := bookCollection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "book not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(book)
}

type CheckInBooksReq struct {
	BookIDs []string `json:"book_ids"`
	UserID  string   `json:"user_id"`
}

func CheckInBooks(c *fiber.Ctx) error {
	if !services.IsKiosk(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	data := new(CheckInBooksReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	userID, userIDErr := bson.ObjectIDFromHex(data.UserID)
	if userIDErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bookCollection := config.GetBookCollection()
	historyCollection := config.GetHistoryCollection()
	for _, bookIDHex := range data.BookIDs {
		bookID, bookIDErr := bson.ObjectIDFromHex(bookIDHex)
		if bookIDErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid book ID: " + bookIDHex,
			})
		}
		var book models.Book
		bookCollection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
		if book.ID.IsZero() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "book not found: " + bookIDHex,
			})
		}
		if book.TakenByUserID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "book is taken: " + bookIDHex,
			})
		}
		_, updateErr := bookCollection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{
			"$set": bson.M{
				"taken_by_user_id": userID,
			},
		})
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check in book: " + bookIDHex,
			})
		}
		newHistory := bson.M{
			"book_id":   bookID,
			"user_id":   userID,
			"issued_at": time.Now(),
		}
		_, insertErr := historyCollection.InsertOne(ctx, newHistory)
		if insertErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create history for book: " + bookIDHex,
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "books checked in successfully",
	})
}

type ReturnBooksReq struct {
	BookIDs []string `json:"book_ids"`
}

func ReturnBooks(c *fiber.Ctx) error {
	if !services.IsKiosk(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}
	data := new(ReturnBooksReq)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bookCollection := config.GetBookCollection()
	historyCollection := config.GetHistoryCollection()
	for _, bookIDHex := range data.BookIDs {
		bookID, bookIDErr := bson.ObjectIDFromHex(bookIDHex)
		if bookIDErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid book ID: " + bookIDHex,
			})
		}
		var book models.Book
		bookCollection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
		if book.ID.IsZero() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "book not found: " + bookIDHex,
			})
		}
		if book.TakenByUserID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "book is not taken by the user: " + bookIDHex,
			})
		}
		_, updateErr := bookCollection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{
			"$set": bson.M{
				"taken_by_user_id": nil,
			},
		})
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to return book: " + bookIDHex,
			})
		}
		currentTime := time.Now()
		_, historyUpdateErr := historyCollection.UpdateOne(ctx, bson.M{
			"book_id": bookID,
			"user_id": *book.TakenByUserID,
			"returned_at": bson.M{
				"$exists": false,
			},
		}, bson.M{
			"$set": bson.M{
				"returned_at": &currentTime,
			},
		})
		if historyUpdateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update history for book: " + bookIDHex,
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "books returned successfully",
	})
}

type GetAllBooksReq struct {
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
	AvailableOnly bool   `json:"available_only"`
	Genre         string `json:"genre"`
	Author        string `json:"author"`
	Publisher     string `json:"publisher"`
	ShelfID       string `json:"shelf_id"`
	Search        string `json:"search"`
}

type BookWithShelf struct {
	models.Book  `bson:",inline"`
	ShelfAddress string `bson:"shelf_address" json:"shelf_address"`
}

func GetAllBooks(c *fiber.Ctx) error {
	if !services.IsAdmin(c) && !services.IsNormalUser(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized access",
		})
	}

	req := new(GetAllBooksReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bookCollection := config.GetBookCollection()

	filter := bson.M{}
	if req.AvailableOnly {
		filter["taken_by_user_id"] = nil
	}
	if req.Genre != "" {
		filter["genre"] = bson.M{"$regex": req.Genre, "$options": "i"}
	}
	if req.Author != "" {
		filter["author"] = bson.M{"$regex": req.Author, "$options": "i"}
	}
	if req.Publisher != "" {
		filter["publisher"] = bson.M{"$regex": req.Publisher, "$options": "i"}
	}
	if req.ShelfID != "" {
		shelfID, err := bson.ObjectIDFromHex(req.ShelfID)
		if err == nil {
			filter["shelf_id"] = shelfID
		}
	}
	if req.Search != "" {
		searchRegex := bson.M{"$regex": req.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"title": searchRegex},
			{"author": searchRegex},
			{"publisher": searchRegex},
			{"isbn": searchRegex},
			{"genre": searchRegex},
		}
	}

	total, err := bookCollection.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to count books",
		})
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "shelves"},
			{Key: "localField", Value: "shelf_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "shelf_details"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$shelf_details"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "shelf_address", Value: "$shelf_details.address"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "shelf_details", Value: 0},
		}}},
	}

	cursor, err := bookCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch books",
		})
	}
	defer cursor.Close(ctx)

	var books []BookWithShelf
	if err = cursor.All(ctx, &books); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode books",
		})
	}

	if books == nil {
		books = []BookWithShelf{}
	}

	lastPage := (total + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": books,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": lastPage,
			"limit":     limit,
		},
	})
}
