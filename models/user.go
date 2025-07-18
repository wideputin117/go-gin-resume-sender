package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	PhoneNumber int64              `bson:"phoneNumber" json:"phoneNumber"`
	Password    int64              `bson:"password" json:"password"`
	Email       string             `bson:"email" json:"email"`
}
