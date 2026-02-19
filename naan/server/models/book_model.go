package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Book struct {
	ID            bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Title         string         `bson:"title" json:"title"`
	Author        string         `bson:"author" json:"author"`
	Publisher     string         `bson:"publisher" json:"publisher"`
	ISBN          string         `bson:"isbn" json:"isbn"`
	Genre         string         `bson:"genre" json:"genre"`
	ShelfID       bson.ObjectID  `bson:"shelf_id" json:"shelf_id"`
	AddedAt       time.Time      `bson:"added_at" json:"added_at"`
	TakenByUserID *bson.ObjectID `bson:"taken_by_user_id,omitempty" json:"taken_by_user_id,omitempty"`
	Row           int            `bson:"row" json:"row"`
	Column        int            `bson:"column" json:"column"`
}

type PublicBook struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string        `bson:"title" json:"title"`
	Author    string        `bson:"author" json:"author"`
	Publisher string        `bson:"publisher" json:"publisher"`
	ISBN      string        `bson:"isbn" json:"isbn"`
	Genre     string        `bson:"genre" json:"genre"`
}

type History struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BookID     bson.ObjectID `bson:"book_id" json:"book_id"`
	UserID     bson.ObjectID `bson:"user_id" json:"user_id"`
	IssuedAt   time.Time     `bson:"issued_at" json:"issued_at"`
	ReturnedAt *time.Time    `bson:"returned_at,omitempty" json:"returned_at,omitempty"`
}
