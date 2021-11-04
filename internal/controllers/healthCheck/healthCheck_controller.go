package healthCheck

import (
	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/apis/healthCheck"
	"net/http"
)

// HealthCheckHandler - returns health status of the API
func HealthCheckHandler(c *gin.Context) {
	healthCheck := healthCheck.HealthCheck{Healthy: true}

	c.JSON(http.StatusOK, healthCheck)
}