package wines

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *WineService) {
	r.POST("/wines", handleCreateWine(s))
	r.GET("/wines/:id", handleGetWine(s))
	r.GET("/wines", handleListWines(s))
	r.PUT("/wines/:id", handleUpdateWine(s))
	r.DELETE("/wines/:id", handleDeleteWine(s))
}

func handleCreateWine(s *WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var wine Wine
		if err := c.BindJSON(&wine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		id, err := s.CreateWine(c, &wine)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id": id,
		})
	}
}

func handleGetWine(s *WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		wine, err := s.GetWine(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, wine)
	}
}

func handleListWines(s *WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wines, err := s.ListWines(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, wines)
	}
}

func handleUpdateWine(s *WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var wine Wine
		if err := c.BindJSON(&wine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := s.UpdateWine(c, id, &wine); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func handleDeleteWine(s *WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.DeleteWine(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
