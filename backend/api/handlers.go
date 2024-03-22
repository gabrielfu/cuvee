package api

import (
	"cuvee/domain/wines"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {}

func handleCreateWine(s *wines.WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wines",
		})
	}
}

func handleGetWine(s *wines.WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wines",
		})
	}
}

func handleListWines(s *wines.WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wines",
		})
	}
}

func handleUpdateWine(s *wines.WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wines",
		})
	}
}

func handleDeleteWine(s *wines.WineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "wines",
		})
	}
}
