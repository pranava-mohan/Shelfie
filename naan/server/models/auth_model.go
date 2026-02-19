package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email" json:"email"`
	Phone     string        `bson:"phone" json:"phone"`
	DauthID   string        `bson:"dauth_id" json:"dauth_id"`
	GoogleID  string        `bson:"google_id" json:"google_id"`
	CreatedAt time.Time     `bson:"createdAt" json:"created_at"`
}

type PublicUser struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string        `bson:"name" json:"name"`
}

type AdminUser struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string        `bson:"username" json:"username"`
	PasswordHash string        `bson:"password_hash" json:"password_hash"`
	CreatedAt    time.Time     `bson:"createdAt" json:"created_at"`
}
