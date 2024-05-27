package vintagecharts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *VintageChartService) {
	r.GET("/vintage_charts", handleListVintageCharts(s))
	r.GET("/vintage_charts/:symbol/regions", handleListRegions(s))
	r.POST("/vintage_charts/:symbol/suggest", handleSuggestRegion(s))
	r.GET("/vintage_charts/:symbol/ratings", handleGetRating(s))
}

func handleListVintageCharts(s *VintageChartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		vcs := s.ListVintageCharts()
		c.JSON(http.StatusOK, vcs)
	}
}

func handleListRegions(s *VintageChartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		regions, err := s.ListRegions(symbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, regions)
	}
}

func handleSuggestRegion(s *VintageChartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")

		var request SuggestRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		provider, err := s.getProvider(symbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": err.Error(),
			})
			return
		}

		region, err := s.SuggestRegion(c, request, provider.ListRegions())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"type":  "internal",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"region": region,
		})
	}
}

func handleGetRating(s *VintageChartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		region := c.Query("region")
		vintage := c.Query("vintage")
		rating, err := s.GetRating(c, symbol, region, vintage)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"type":  "not found",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, rating)
	}
}
