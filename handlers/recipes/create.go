package recipes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/InveterateCoder/recipes-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (h *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.PublishedAt = time.Now()
	ctx := c.Request.Context()
	res, err := h.collection.InsertOne(ctx, recipe)
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

	h.redisClient.Del(ctx, "recipes")
	c.JSON(http.StatusCreated, recipe)
}
