package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Established int32  `bson:"established" json:"established,omitempty"`
}

