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
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

// swagger:model
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/Recipe'
func ListRecipesHandler(c *gin.Context) {
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
// - name: recipe
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/Recipe'
// responses:
//   '200':
//     description: Successful operation
//     schema:
//       $ref: '#/definitions/Recipe'
//   '400':
//     description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Updates a recipe by ID.
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   required: true
//   type: string
// - name: recipe
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/Recipe'
// responses:
//   '200':
//     description: Successful operation
//     schema:
//       $ref: '#/definitions/Recipe'
//   '400':
//     description: Invalid input
//   '404':
//     description: Recipe not found
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipe.ID = id
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Deletes a recipe by ID.
// ---
// parameters:
// - name: id
//   in: path
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successful operation
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//   '404':
//     description: Recipe not found
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
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
// - name: tag
//   in: query
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successful operation
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/Recipe'
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func init() {
	if file, err := os.ReadFile("recipes.json"); err != nil {
		panic(err)
	} else {
		if err = json.Unmarshal(file, &recipes); err != nil {
			panic(err)
		}
	}
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
