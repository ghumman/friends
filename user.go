package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User ...
type User struct {
	ID         primitive.ObjectID `bson:"_id, omitempty"`
	FirstName  string             `bson:"firstName" json:"firstName"`
	LastName   string             `bson:"lastName" json:"lastName"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	Salt       string             `bson:"salt" json:"salt"`
	Password   string             `bson:"password" json:"password"`
	Email      string             `bson:"email" json:"email"`
	AuthType   int                `bson:"authType" json:"authType"`
	Token      string             `bson:"token" json:"token"`
	ResetToken string             `bson:"resetToken" json:"resetToken"`
}
