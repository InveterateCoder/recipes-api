package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// swagger:model
type Recipe struct {
	// swagger:ignore
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Tags         []string      `json:"tags" bson:"tags"`
	Ingredients  []string      `json:"ingredients" bson:"ingredients"`
	Instructions []string      `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time     `json:"publishedAt" bson:"publishedAt"`
}
