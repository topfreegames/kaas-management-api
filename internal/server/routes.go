package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	"github.com/topfreegames/kaas-management-api/api/healthCheck"
	nodegroupv1 "github.com/topfreegames/kaas-management-api/api/nodeGroup/v1"
	"github.com/topfreegames/kaas-management-api/internal/controller"
	"net/http"
)

type RouterConfig struct {
	controller controller.ControllerConfig
	router     *gin.Engine
}

// path returns a path with trailing slash
func path(path string) string {
	return path + "/"
}

// param returns a path parameter in the gin notation as string
func param(param string) string {
	return ":" + param + "/"
}

func (r RouterConfig) setupRoutes() {
	r.setupClusterV1Routes()
	r.setupHealthCheckRoutes()
	r.setupDocsRoutes()
}

func (r RouterConfig) setupClusterV1Routes() {
	r.router.Handle(http.MethodGet, clusterv1.Endpoint.Path, r.controller.ClusterListHandler)
	r.router.Handle(http.MethodGet, clusterv1.Endpoint.Path+param(clusterv1.ClusterNameParameter), r.controller.ClusterHandler)
	r.router.Handle(http.MethodGet, clusterv1.Endpoint.Path+param(clusterv1.ClusterNameParameter)+path(nodegroupv1.Endpoint.EndpointName), r.controller.NodeGroupListByClusterHandler)
	r.router.Handle(http.MethodGet, clusterv1.Endpoint.Path+param(clusterv1.ClusterNameParameter)+path(nodegroupv1.Endpoint.EndpointName)+param(nodegroupv1.NodeGroupNameParameter), r.controller.NodeGroupByClusterHandler)
}

func (r RouterConfig) setupHealthCheckRoutes() {
	r.router.Handle(http.MethodGet, healthCheck.Endpoint.Path, controller.HealthCheckHandler)
}

func (r RouterConfig) setupDocsRoutes() {
	r.router.GET("/docs", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/docs/swagger/index.html")
	})
	r.router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
