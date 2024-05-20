package regions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *RegionService) {
	r.GET("/regions/wines/:wineId", handleListRegions(s))
	r.GET("/regions/wines/:wineId/vcs/:vcSymbol", handleGetRegion(s))
	r.POST("/regions/wines/:wineId", handleCreateRegion(s))
	r.PUT("/regions/wines/:wineId/vcs/:vcSymbol", handleUpdateRegion(s))
	r.DELETE("/regions/wines/:wineId/vcs/:vcSymbol", handleDeleteRegion(s))
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
		vcSymbol := c.Param("vcSymbol")
		region, err := s.GetRegion(c, wineID, vcSymbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vc %s not found", wineID, vcSymbol),
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
		if _, err := s.GetRegion(c, wineID, region.VCSymbol); err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"type":  "conflict",
				"error": fmt.Sprintf("region for wine %s and vc %s already exists", wineID, region.VCSymbol),
			})
			return
		}

		if err := s.CreateRegion(c, &Region{
			WineID:   wineID,
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

func handleUpdateRegion(s *RegionService) gin.HandlerFunc {
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

		if _, err := s.GetRegion(c, wineID, region.VCSymbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vc %s not found", wineID, region.VCSymbol),
			})
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

func handleDeleteRegion(s *RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wineID := c.Param("wineId")
		vcSymbol := c.Param("vcSymbol")

		if _, err := s.GetRegion(c, wineID, vcSymbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": fmt.Sprintf("region for wine %s and vc %s not found", wineID, vcSymbol),
			})
		}

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
