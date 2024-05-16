package ratings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *RatingService) {
	r.GET("/ratings/regions/wines/:wineId", handleListRegions(s))
	r.GET("/ratings/regions/wines/:wineId/vc/:vcSymbol", handleGetRegion(s))
	r.POST("/ratings/regions/wines/:wineId", handleCreateRegion(s))
	r.PUT("/ratings/regions/wines/:wineId/vc/:vcSymbol", handleUpdateRegion(s))
	r.DELETE("/ratings/regions/wines/:wineId/vc/:vcSymbol", handleDeleteRegion(s))
}

func handleListRegions(s *RatingService) gin.HandlerFunc {
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

func handleGetRegion(s *RatingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		vcSymbol := c.Param("vcSymbol")
		region, err := s.GetRegion(c, wineID, vcSymbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, region)
	}
}

func handleCreateRegion(s *RatingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var region Region
		if err := c.BindJSON(&region); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"type":  "bad request",
				"error": err.Error(),
			})
			return
		}

		if err := s.CreateRegion(c, &Region{
			WineID:   region.WineID,
			VCSymbol: region.VCSymbol,
			Region:   region.Region,
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

func handleUpdateRegion(s *RatingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		vcSymbol := c.Param("vcSymbol")

		var region Region
		if err := c.BindJSON(&region); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"type":  "bad request",
				"error": err.Error(),
			})
			return
		}

		if err := s.UpdateRegion(c, &Region{
			WineID:   wineID,
			VCSymbol: vcSymbol,
			Region:   region.Region,
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

func handleDeleteRegion(s *RatingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		vcSymbol := c.Param("vcSymbol")
		if err := s.DeleteRegion(c, wineID, vcSymbol); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}