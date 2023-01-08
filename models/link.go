package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	ShortURL  string             `bson:"short_url" json:"short_url"`
	LongURL   string             `bson:"long_url" json:"long_url"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
