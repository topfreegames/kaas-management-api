package server

import (
	"github.com/gin-gonic/gin"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	"github.com/topfreegames/kaas-management-api/api/healthCheck"
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

func (r RouterConfig) setupClusterV1Routes() {
	clusterv1.Endpoint.CreatePrivateRouterGroup(r.router)
	clusterv1.Endpoint.CreateRoute("GET", "/", r.controller.ClusterListHandler)
	clusterv1.Endpoint.CreateRoute("GET", "/:" + clusterv1.ClusterNameParameter, r.controller.ClusterHandler)
}

func (r RouterConfig) setupHealthCheckRoutes() {
	healthCheck.Endpoint.CreatePublicRouterGroup(r.router)
	healthCheck.Endpoint.CreateRoute("GET", "/", controller.HealthCheckHandler)
}
