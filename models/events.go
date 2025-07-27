package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Location    string             `json:"location" bson:"location"`
	Date        string             `json:"date" bson:"date"`
	CreatedBy   string             `json:"createdBy" bson:"createdBy"` // <-- Add this
}