package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	// 1. Define command-line flags
	username := flag.String("username", "", "Username for the new admin user (required)")
	password := flag.String("password", "", "Password for the new admin user (required)")

	flag.Parse()

	// 2. Validate input
	if *username == "" || *password == "" {
		log.Println("Error: --username and --password flags are required.")
		flag.Usage()
		return
	}

	config.LoadEnv()
	config.ConnectDB()

	// 4. Check if user already exists
	var existingUser models.AdminUser
	userCollection := config.GetAdminUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"username": *username}).Decode(&existingUser)

	if err != mongo.ErrNoDocuments {
		log.Fatalf("error: %v", err)
		log.Fatalf("User already exists prolly")
	}

	hashedPassword, err := hashPassword(*password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	newUser := models.AdminUser{
		Username:     *username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	_, insertErr := userCollection.InsertOne(ctx, newUser)
	if insertErr != nil {
		log.Fatalf("Failed to create admin user: %v", insertErr)
	}

	fmt.Printf("✅ Successfully created user!\n")
	fmt.Printf("   Username: %s\n", newUser.Username)
}
