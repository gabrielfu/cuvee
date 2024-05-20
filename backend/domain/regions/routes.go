package regions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *RegionService) {
	r.GET("/regions/wines/:wineId", handleListRegions(s))
	r.GET("/regions/wines/:wineId/vintage_charts/:symbol", handleGetRegion(s))
	r.POST("/regions/wines/:wineId", handleCreateRegion(s))
	r.PUT("/regions/wines/:wineId/vintage_charts/:symbol", handleUpdateRegion(s))
	r.DELETE("/regions/wines/:wineId/vintage_charts/:symbol", handleDeleteRegion(s))
}

func handleListRegions(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		regions, err := s.ListRegions(c, wineID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, regions)
	}
}

func handleGetRegion(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		symbol := c.Param("symbol")
		region, err := s.GetRegion(c, wineID, symbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vintage chart %s not found", wineID, symbol),
			})
			return
		}
		c.JSON(http.StatusOK, region)
	}
}

func handleCreateRegion(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var region Region
		if err := c.BindJSON(&region); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"type":  "bad request",
				"error": err.Error(),
			})
			return
		}

		wineID := c.Param("wineId")
		if _, err := s.GetRegion(c, wineID, region.Symbol); err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"type":  "conflict",
				"error": fmt.Sprintf("region for wine %s and vintage chart %s already exists", wineID, region.Symbol),
			})
			return
		}

		if err := s.CreateRegion(c, &Region{
			WineID: wineID,
			Symbol: region.Symbol,
			Region: region.Region,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{})
	}
}

func handleUpdateRegion(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		symbol := c.Param("symbol")

		var region Region
		if err := c.BindJSON(&region); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"type":  "bad request",
				"error": err.Error(),
			})
			return
		}

		if _, err := s.GetRegion(c, wineID, region.Symbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vintage chart %s not found", wineID, region.Symbol),
			})
		}

		if err := s.UpdateRegion(c, &Region{
			WineID: wineID,
			Symbol: symbol,
			Region: region.Region,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func handleDeleteRegion(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		symbol := c.Param("symbol")

		if _, err := s.GetRegion(c, wineID, symbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vintage chart %s not found", wineID, symbol),
			})
		}

		if err := s.DeleteRegion(c, wineID, symbol); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
