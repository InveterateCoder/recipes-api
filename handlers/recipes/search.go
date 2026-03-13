package recipes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/InveterateCoder/recipes-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (h *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := strings.TrimSpace(c.Query("tag"))
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tag query parameter is required",
		})
		return
	}

	ctx := c.Request.Context()
	filter := bson.M{
		"tags": bson.Regex{
			Pattern: "^" + regexp.QuoteMeta(tag) + "$",
			Options: "i",
		},
	}

	cur, err := h.collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	listOfRecipes := make([]models.Recipe, 0)
	for cur.Next(ctx) {
		var recipe models.Recipe
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
