package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookData struct {
	Title  string
	Author string
	Genre  string
}

var realBooks = []BookData{
	{"To Kill a Mockingbird", "Harper Lee", "Fiction"},
	{"1984", "George Orwell", "Science Fiction"},
	{"The Great Gatsby", "F. Scott Fitzgerald", "Fiction"},
	{"Pride and Prejudice", "Jane Austen", "Romance"},
	{"The Catcher in the Rye", "J.D. Salinger", "Fiction"},
	{"The Hobbit", "J.R.R. Tolkien", "Fantasy"},
	{"Fahrenheit 451", "Ray Bradbury", "Science Fiction"},
	{"The Lord of the Rings", "J.R.R. Tolkien", "Fantasy"},
	{"Animal Farm", "George Orwell", "Satire"},
	{"The Diary of a Young Girl", "Anne Frank", "Biography"},
	{"The Alchemist", "Paulo Coelho", "Fiction"},
	{"Harry Potter and the Sorcerer's Stone", "J.K. Rowling", "Fantasy"},
	{"The Little Prince", "Antoine de Saint-Exupéry", "Fiction"},
	{"A Game of Thrones", "George R.R. Martin", "Fantasy"},
	{"The Da Vinci Code", "Dan Brown", "Thriller"},
	{"The Hunger Games", "Suzanne Collins", "Science Fiction"},
	{"The Kite Runner", "Khaled Hosseini", "Fiction"},
	{"Life of Pi", "Yann Martel", "Fiction"},
	{"The Book Thief", "Markus Zusak", "Historical Fiction"},
	{"Gone Engineer", "Pranava Mohan", "Biography"},
	{"Sapiens: A Brief History of Humankind", "Yuval Noah Harari", "Non-Fiction"},
	{"Becoming", "Michelle Obama", "Biography"},
	{"Educated", "Tara Westover", "Biography"},
	{"The Immortal Life of Henrietta Lacks", "Rebecca Skloot", "Non-Fiction"},
	{"Steve Jobs", "Walter Isaacson", "Biography"},
	{"Elon Musk", "Walter Isaacson", "Biography"},
	{"Thinking, Fast and Slow", "Daniel Kahneman", "Psychology"},
	{"Atomic Habits", "James Clear", "Self-Help"},
	{"Deep Work", "Cal Newport", "Self-Help"},
	{"Clean Code", "Robert C. Martin", "Technology"},
	{"The Pragmatic Programmer", "Andrew Hunt & David Thomas", "Technology"},
	{"Introduction to Algorithms", "Thomas H. Cormen", "Technology"},
	{"Design Patterns", "Erich Gamma", "Technology"},
	{"Cracking the Coding Interview", "Gayle Laakmann McDowell", "Technology"},
	{"Dune", "Frank Herbert", "Science Fiction"},
	{"Neuromancer", "William Gibson", "Science Fiction"},
	{"Snow Crash", "Neal Stephenson", "Science Fiction"},
	{"Foundation", "Isaac Asimov", "Science Fiction"},
	{"Brave New World", "Aldous Huxley", "Science Fiction"},
	{"The Hitchhiker's Guide to the Galaxy", "Douglas Adams", "Science Fiction"},
}

func main() {

	config.LoadEnv()
	config.ConnectDB()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shelfColl := config.GetShelfCollection()
	bookColl := config.GetBookCollection()

	var shelves []models.Shelf
	cursor, err := shelfColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to fetch shelves: %v", err)
	}
	if err = cursor.All(ctx, &shelves); err != nil {
		log.Fatalf("Failed to decode shelves: %v", err)
	}

	if len(shelves) == 0 {
		log.Fatalf("No shelves found. Please create at least one shelf before seeding books.")
	}

	fmt.Printf("Found %d shelves. Starting seed...\n", len(shelves))

	booksToInsert := make([]interface{}, 0, 100)
	for i := 0; i < 100; i++ {
		randomShelf := shelves[rand.Intn(len(shelves))]
		randomBook := realBooks[rand.Intn(len(realBooks))]

		book := models.Book{
			Title:     randomBook.Title,
			Author:    randomBook.Author,
			Publisher: "Generic Publisher",
			ISBN:      fmt.Sprintf("978-%d-%d-%d-%d", rand.Intn(10), rand.Intn(1000), rand.Intn(100), rand.Intn(10)),
			Genre:     randomBook.Genre,
			ShelfID:   randomShelf.ID,
			AddedAt:   time.Now(),
			Row:       rand.Intn(5) + 1,
			Column:    rand.Intn(5) + 1,
		}
		booksToInsert = append(booksToInsert, book)
	}

	result, err := bookColl.InsertMany(ctx, booksToInsert)
	if err != nil {
		log.Fatalf("Failed to seed books: %v", err)
	}

	fmt.Printf("✅ Successfully seeded %d books with realistic data!\n", len(result.InsertedIDs))
}
