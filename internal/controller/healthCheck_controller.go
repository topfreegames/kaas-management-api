package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/api/healthCheck"
)

// HealthCheckHandler - returns health status of the API
func HealthCheckHandler(c *gin.Context) {
	healthCheck := healthCheck.HealthCheck{Healthy: true}

	c.JSON(http.StatusOK, healthCheck)
}
