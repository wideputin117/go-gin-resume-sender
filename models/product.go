package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID primitive.ObjectID   `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Category primitive.ObjectID `bson:"category" json:"category"`
	Stock int64 `bson:"stock" json:"stock"`
	Price int16 `bson:"price" json:"price"`
}