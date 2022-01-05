package server

import (
	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/internal/controller"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
)

// InitServer - Initializes the serves
func InitServer(k8sInstance *k8s.Kubernetes) error {

	router := gin.Default()
	controllerInstance := controller.ConfigureControllers(k8sInstance)

	routerConfig := &RouterConfig{
		controller: controllerInstance,
		router:     router,
	}
	routerConfig.setupRoutes()
	return router.Run()
}
