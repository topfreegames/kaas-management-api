package server

import (
	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/docs"
	"github.com/topfreegames/kaas-management-api/internal/controller"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
)

// @securityDefinitions.basic  BasicAuth

// InitServer - Initializes the serves
func InitServer(k8sInstance *k8s.Kubernetes) error {

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Kubernetes as a service API"
	docs.SwaggerInfo.Description = "K8s Clusters management API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}

	router := gin.Default()
	controllerInstance := controller.ConfigureControllers(k8sInstance)

	routerConfig := &RouterConfig{
		controller: controllerInstance,
		router:     router,
	}
	routerConfig.setupRoutes()
	return router.Run()
}
