package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	healthCheckv1 "github.com/topfreegames/kaas-management-api/api/healthCheck"
	"github.com/topfreegames/kaas-management-api/internal/controller"
)

type RouterConfig struct {
	controller controller.ControllerConfig
	router     *gin.Engine
}

func (r RouterConfig) setupRoutes() {
	r.setupClusterV1Routes()
	r.setupHealthCheckRoutes()
}

func (r RouterConfig) setupHealthCheckRoutes() {
	r.router.Handle("GET", fmt.Sprintf("/%s", healthCheckv1.Endpoint), controller.HealthCheckHandler)
}

func (r RouterConfig) setupClusterV1Routes() {

	group := r.router.Group(fmt.Sprintf("%s/%s", clusterv1.Version, clusterv1.Endpoint))
	group.Handle("GET", "/", r.controller.ClusterListHandler)
	group.Handle("GET", "/:name", r.controller.ClusterHandler)
}
