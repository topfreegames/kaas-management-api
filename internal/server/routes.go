package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	clusterv1 "github.com/topfreegames/kaas-management-api/apis/cluster/v1"
	healthCheckv1 "github.com/topfreegames/kaas-management-api/apis/healthCheck"
	"github.com/topfreegames/kaas-management-api/internal/controllers/cluster"
	"github.com/topfreegames/kaas-management-api/internal/controllers/healthCheck"
)

func setupHealthCheckRoutes(router *gin.Engine) {
	router.Handle("GET", fmt.Sprintf("/%s", healthCheckv1.Endpoint), healthCheck.HealthCheckHandler)
}

func setupClusterV1Routes(router *gin.Engine) {
	group := router.Group(fmt.Sprintf("%s/%s", clusterv1.Version, clusterv1.Endpoint))
	group.Handle("GET", "/", cluster.ClusterListHandler)
	group.Handle("GET", "/:name", cluster.ClusterHandler)
}
