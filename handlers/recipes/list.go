package recipes

import (
	"net/http"

	"github.com/InveterateCoder/recipes-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (h *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	cur, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)
	recipes := make([]models.Recipe, 0)
	for cur.Next(ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}
