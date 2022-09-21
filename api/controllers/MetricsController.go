package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//PublicGetMetrics is a public endpoint to get metrics
func PublicGetMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"views": 0,
	})
}

//PublicLogMetric is a public endpoint to log metrics
func PublicLogMetric(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"views": 0,
	})
}
