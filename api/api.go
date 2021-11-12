package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
)

type ApiEndpoint struct {
    version string
    endpoint string
    routerGroup *gin.RouterGroup
}

// NewApiEndpoint Creates a new ApiEndpoint structure representing an versioned API endpoint
func NewApiEndpoint(version string, endpoint string) *ApiEndpoint {
    if endpoint == "" {
        log.Fatalf("Endpoint pattern can't be empty")
    }
    return &ApiEndpoint{
        version:  version,
        endpoint: endpoint,
    }
}

// CreatePrivateRouterGroup creates a new Gin Router group for the endpoint with an authentication middleware
func (a *ApiEndpoint) CreatePrivateRouterGroup(engine *gin.Engine) {
    a.routerGroup = engine.Group(a.getEndpoint())
    // TODO add authentication middleware
    // a.routerGroup.Use()
}

// CreatePublicRouterGroup creates a new Gin Router group for the endpoint publicly exposed
func (a *ApiEndpoint) CreatePublicRouterGroup(engine *gin.Engine) {
    a.routerGroup = engine.Group(a.getEndpoint())
}

// CreateRouterGroup creates a new route in the Gin router group of the endpoint
func (a *ApiEndpoint) CreateRoute(method string, pattern string, handlerFunc gin.HandlerFunc) {
    a.routerGroup.Handle(method, pattern, handlerFunc)
}

// GetEndpoint Returns the string pattern of the endpoint with the version
func (a *ApiEndpoint) getEndpoint() string {
    if a.version != "" {
        return fmt.Sprintf("/%s/%s", a.version, a.endpoint)
    }
    return fmt.Sprintf("/%s", a.endpoint)
}