package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/InveterateCoder/recipes-api/handlers/recipes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var client *mongo.Client
var db *mongo.Database
var recipesHandler *recipes.RecipesHandler

func init() {
	godotenv.Load()

	var (
		err      error
		dbName   = os.Getenv("MONGO_DB")
		mongoUri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=%s",
			os.Getenv("MONGO_USER"),
			os.Getenv("MONGO_PASSWORD"),
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT"),
			dbName,
			os.Getenv("MONGO_AUTH_SOURCE"),
		)
	)
	if client, err = mongo.Connect(options.Client().ApplyURI(mongoUri)); err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB")
	db = client.Database(dbName)
	recipesHandler = recipes.NewRecipeHandler(db.Collection("recipes"))
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.Run()
}
