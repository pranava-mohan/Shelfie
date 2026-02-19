package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Shelf struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Address string        `bson:"address" json:"address"`
}
