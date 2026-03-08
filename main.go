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
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Successful operation
//	  schema:
//	    type: array
//	    items:
//	      $ref: '#/definitions/Recipe'
func ListRecipesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	collection := db.Collection("recipes")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)
	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation POST /recipes recipes createRecipe
// Creates a recipe.
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: recipe
//     in: body
//     required: true
//     schema:
//     $ref: '#/definitions/Recipe'
//
// responses:
//
//	'200':
//	  description: Successful operation
//	  schema:
//	    $ref: '#/definitions/Recipe'
//	'400':
//	  description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.PublishedAt = time.Now()
	ctx := c.Request.Context()
	collection := db.Collection("recipes")
	res, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	id, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Inserted id is not an ObjectID"})
		return
	}
	recipe.ID = id
	c.JSON(http.StatusCreated, recipe)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Updates a recipe by ID.
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     required: true
//     type: string
//   - name: recipe
//     in: body
//     required: true
//     schema:
//     $ref: '#/definitions/Recipe'
//
// responses:
//
//	'200':
//	  description: Successful operation
//	  schema:
//	    $ref: '#/definitions/Recipe'
//	'400':
//	  description: Invalid input
//	'404':
//	  description: Recipe not found
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	collection := db.Collection("recipes")
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: recipe.Name},
				{Key: "instructions", Value: recipe.Instructions},
				{Key: "ingredients", Value: recipe.Ingredients},
				{Key: "tags", Value: recipe.Tags},
			},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Deletes a recipe by ID.
// ---
// parameters:
//   - name: id
//     in: path
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Successful operation
//	  schema:
//	    type: object
//	    properties:
//	      message:
//	        type: string
//	'404':
//	  description: Recipe not found
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	collection := db.Collection("recipes")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Searches recipes by tag.
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Successful operation
//	  schema:
//	    type: array
//	    items:
//	      $ref: '#/definitions/Recipe'
func SearchRecipesHandler(c *gin.Context) {
	tag := strings.TrimSpace(c.Query("tag"))
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tag query parameter is required",
		})
		return
	}

	ctx := c.Request.Context()
	collection := db.Collection("recipes")
	filter := bson.M{
		"tags": bson.Regex{
			Pattern: "^" + regexp.QuoteMeta(tag) + "$",
			Options: "i",
		},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	listOfRecipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		if err := cur.Decode(&recipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		listOfRecipes = append(listOfRecipes, recipe)
	}
	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

var client *mongo.Client
var db *mongo.Database

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
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/recipes", ListRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
