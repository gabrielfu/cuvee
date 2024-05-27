package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *ImageService) {
	r.POST("/images/search", handleImageSearch(s))
}

func handleImageSearch(s *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ImageSearchRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if request.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
			return
		}

		resp, err := s.Search(c, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, &resp)
	}
}
