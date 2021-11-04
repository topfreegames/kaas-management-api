package server

import (
    "github.com/gin-gonic/gin"
)

// InitServer - Initializes the serves
func InitServer() error {
    router := gin.Default()
    setupClusterV1Routes(router)
    setupHealthCheckRoutes(router)
    return router.Run()
}
