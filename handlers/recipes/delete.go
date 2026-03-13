package recipes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (h *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	res, err := h.collection.DeleteOne(ctx, bson.M{"_id": objectId})
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
