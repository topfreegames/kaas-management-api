package server

import (
    "github.com/gin-gonic/gin"
    "github.com/topfreegames/kaas-management-api/internal/controllers/cluster"
)

// InitServer - Initializes the serves
func InitServer() error {
    router := gin.Default()
    router.GET("/ping", cluster.PingHandler)
    router.GET("/cluster", cluster.ClusterHandler)
    return router.Run()
}
