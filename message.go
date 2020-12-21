package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message ...
type Message struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Message     string             `bson:"message" json:"message"`
	MessageFrom User               `bson:"messageFrom" json:"messageFrom"`
	MessageTo   User               `bson:"messageTo" json:"messageTo"`
	SentAt      time.Time          `bson:"sentAt" json:"sentAt"`
}
