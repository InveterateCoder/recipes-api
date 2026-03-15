package recipes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/InveterateCoder/recipes-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (h *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
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
	_, err = h.collection.UpdateOne(ctx, bson.M{
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
	h.redisClient.Del(ctx, "recipes")
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}
