package recipes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/InveterateCoder/recipes-api/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	recipes := make([]models.Recipe, 0)
	val, cacheerr := h.redisClient.Get(ctx, "recipes").Result()
	switch cacheerr {
	case nil:
		log.Printf("Requesting to Redis")
		if err := json.Unmarshal([]byte(val), &recipes); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case redis.Nil:
		log.Printf("Requesting to MongoDB")
		cur, err := h.collection.Find(ctx, bson.M{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		if err := cur.Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if data, err := json.Marshal(recipes); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			if err := h.redisClient.Set(ctx, "recipes", string(data), 0).Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": cacheerr.Error()})
	}
	c.JSON(http.StatusOK, recipes)
}
