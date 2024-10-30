package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description"`
	Title       string             `json:"title" bson:"title"`
	Completed   bool               `json:"completed" bson:"completed"`
}
