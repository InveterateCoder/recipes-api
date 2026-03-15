// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
// Schema: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Arthur Grigoryan <inveterate.coder@gmail.com> https://inveteratecoder.github.io
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package recipes

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

func NewRecipeHandler(collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		redisClient: redisClient,
	}
}
