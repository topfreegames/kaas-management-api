package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type ApiEndpoint struct {
	version      string
	endpoint     string
	RouterGroup  *gin.RouterGroup
	Router       *gin.Engine
	EndpointPath string
}

// NewApiEndpoint Creates a new ApiEndpoint structure representing an versioned API endpoint
func NewApiEndpoint(version string, endpoint string) *ApiEndpoint {
	if endpoint == "" {
		log.Fatalf("Endpoint pattern can't be empty")
	}

	apiEndpoint := &ApiEndpoint{
		version:  version,
		endpoint: endpoint,
	}
	fullPath := apiEndpoint.GetEndpointPath()
	apiEndpoint.EndpointPath = fullPath

	return apiEndpoint
}

// CreatePrivateRouterGroup creates a new Gin Router group for the endpoint with an authentication middleware
func (a *ApiEndpoint) CreatePrivateRouterGroup(engine *gin.Engine) {
	if a.RouterGroup != nil {
		log.Fatalf("Could not configure router for endpoint %s at version: %s: Router group already exists!", a.endpoint, a.version)
	}
	a.Router = engine
	a.RouterGroup = engine.Group(a.GetEndpointPath())
	// TODO add authentication middleware
	// a.RouterGroup.Use()
}

// CreatePublicRouterGroup creates a new Gin Router group for the endpoint publicly exposed
func (a *ApiEndpoint) CreatePublicRouterGroup(engine *gin.Engine) {
	if a.RouterGroup != nil {
		log.Fatalf("Could not configure router for endpoint %s at version: %s: Router group already exists!", a.endpoint, a.version)
	}
	a.Router = engine
	a.RouterGroup = engine.Group(a.GetEndpointPath())
}

// CreateRouterGroup creates a new route in the Gin router group of the endpoint
func (a *ApiEndpoint) CreateRoute(method string, pattern string, handlerFunc gin.HandlerFunc) {
	a.RouterGroup.Handle(method, pattern, handlerFunc)
}

// GetEndpoint Returns the string pattern of the endpoint with the version
func (a *ApiEndpoint) GetEndpointPath() string {
	if a.version != "" {
		return fmt.Sprintf("/%s/%s", a.version, a.endpoint)
	}
	return fmt.Sprintf("/%s", a.endpoint)
}
