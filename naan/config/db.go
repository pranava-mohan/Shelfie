package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	mongoURI := Env("MONGO_URI", "mongodb://mongodb:27017")
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	DB = client
	EnsureCollectionsAndIndexes(DB.Database(Env("DB_NAME", "my_db")), Collections)
}

type IndexConfig struct {
	Name  string
	Model mongo.IndexModel
}

type CollectionConfig struct {
	Name    string
	Indexes []IndexConfig
}

var Collections = []CollectionConfig{
	{
		Name: "users",
		Indexes: []IndexConfig{
			{
				Name: "email_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true).SetSparse(true).SetName("email_1"),
				},
			},
		},
	},
	{
		Name: "admin_users",
		Indexes: []IndexConfig{
			{
				Name: "username_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "username", Value: 1}},
					Options: options.Index().SetUnique(true).SetSparse(true).SetName("username_1"),
				},
			},
		},
	},
	{
		Name: "books",
		Indexes: []IndexConfig{
			{
				Name: "shelf_id_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "shelf_id", Value: 1}},
					Options: options.Index().SetName("shelf_id_1"),
				},
			},
			{
				Name: "taken_by_user_id_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "taken_by_user_id", Value: 1}},
					Options: options.Index().SetSparse(true).SetName("taken_by_user_id_1"),
				},
			},
			{
				Name: "isbn_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "isbn", Value: 1}},
					Options: options.Index().SetName("isbn_1"),
				},
			},
		},
	},
	{
		Name: "history",
		Indexes: []IndexConfig{
			{
				Name: "book_id_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "book_id", Value: 1}},
					Options: options.Index().SetName("book_id_1"),
				},
			},
			{
				Name: "taken_by_user_id_2",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "taken_by_user_id", Value: 1}},
					Options: options.Index().SetName("taken_by_user_id_2"),
				},
			},
		},
	},
	{
		Name: "kiosks",
		Indexes: []IndexConfig{
			{
				Name: "kiosk_name_1",
				Model: mongo.IndexModel{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetSparse(true).SetUnique(true).SetName("kiosk_name_1"),
				},
			},
		},
	},
}

func EnsureCollectionsAndIndexes(db *mongo.Database, collections []CollectionConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	existingCols, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatalf("❌ failed to list collections: %v", err)
	}

	existingSet := make(map[string]bool)
	for _, name := range existingCols {
		existingSet[name] = true
	}

	for _, cfg := range collections {
		if !existingSet[cfg.Name] {
			if err := db.CreateCollection(ctx, cfg.Name); err != nil {
				log.Printf("⚠️ failed to create collection %s: %v", cfg.Name, err)
			} else {
				log.Printf("✅ created collection: %s", cfg.Name)
			}
		}

		col := db.Collection(cfg.Name)

		cursor, err := col.Indexes().List(ctx)
		if err != nil {
			log.Printf("⚠️ failed to list indexes for %s: %v", cfg.Name, err)
			continue
		}

		existingIdx := map[string]bool{}
		for cursor.Next(ctx) {
			var idx bson.M
			if err := cursor.Decode(&idx); err == nil {
				if name, ok := idx["name"].(string); ok {
					existingIdx[name] = true
				}
			}
		}

		toCreate := []mongo.IndexModel{}
		for _, ic := range cfg.Indexes {
			if !existingIdx[ic.Name] {
				toCreate = append(toCreate, ic.Model)
			}
		}

		if len(toCreate) > 0 {
			_, err = col.Indexes().CreateMany(ctx, toCreate)
			if err != nil {
				log.Printf("⚠️ failed to create indexes on %s: %v", cfg.Name, err)
			} else {
				log.Printf("✅ created %d indexes on %s", len(toCreate), cfg.Name)
			}
		} else {
			log.Printf("ℹ️ all indexes already exist for %s", cfg.Name)
		}
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	dbName := Env("DB_NAME", "my_db")
	collection := DB.Database(dbName).Collection(collectionName)
	return collection
}

func GetUserCollection() *mongo.Collection {
	return GetCollection("users")
}

func GetAdminUserCollection() *mongo.Collection {
	return GetCollection("admin_users")
}

func GetBookCollection() *mongo.Collection {
	return GetCollection("books")
}

func GetShelfCollection() *mongo.Collection {
	return GetCollection("shelves")
}

func GetHistoryCollection() *mongo.Collection {
	return GetCollection("history")
}

func GetKioskCollection() *mongo.Collection {
	return GetCollection("kiosks")
}
