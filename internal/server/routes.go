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

func (r RouterConfig) setupRoutes() {
	r.setupClusterV1Routes()
	r.setupHealthCheckRoutes()
	r.setupDocsRoutes()
}

func (r RouterConfig) setupClusterV1Routes() {
	clusterv1.Endpoint.CreatePrivateRouterGroup(r.router)
	clusterv1.Endpoint.CreateRoute(http.MethodGet, "/", r.controller.ClusterListHandler)
	clusterv1.Endpoint.CreateRoute(http.MethodGet, "/:"+clusterv1.ClusterNameParameter, r.controller.ClusterHandler)
	clusterv1.Endpoint.CreateRoute(http.MethodGet, "/:"+clusterv1.ClusterNameParameter+"/nodegroups", r.controller.NodeGroupListByClusterHandler)
	clusterv1.Endpoint.CreateRoute(http.MethodGet, "/:"+clusterv1.ClusterNameParameter+"/nodegroups"+"/:"+nodegroupv1.NodeGroupNameParameter, r.controller.NodeGroupByClusterHandler)
}

func (r RouterConfig) setupHealthCheckRoutes() {
	healthCheck.Endpoint.CreatePublicRouterGroup(r.router)
	healthCheck.Endpoint.CreateRoute(http.MethodGet, "/", controller.HealthCheckHandler)
}

func (r RouterConfig) setupDocsRoutes() {
	r.router.GET("/docs", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/docs/swagger/index.html")
	})
	r.router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
