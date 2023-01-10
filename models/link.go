package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ShortURL   string             `bson:"short_url" json:"short_url"`
	LongURL    string             `bson:"long_url" json:"long_url"`
	AddedBy    string             `bson:"added_by" json:"added_by"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	AccessedAt time.Time          `bson:"accessed_at" json:"accessed_at"`
}
