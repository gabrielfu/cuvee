package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *ImageService) {
	r.GET("/images/search", handleImageSearch(s))
}

func handleImageSearch(s *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := ImageSearchRequest{
			name:    c.Query("name"),
			vintage: c.Query("vintage"),
			country: c.Query("country"),
			region:  c.Query("region"),
		}
		resp, err := s.Search(c, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
