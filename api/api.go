package api

import (
	"fmt"
	"log"
)

type ApiEndpoint struct {
	// Version Endpoint version
	Version string
	// EndpointNamePath Endpoint name without version
	EndpointName string
	// Path Full path of the endpoint with version and name
	Path string
}

// NewApiEndpoint Creates a new ApiEndpoint structure representing an versioned API EndpointNamePath
func NewApiEndpoint(version string, endpoint string) *ApiEndpoint {
	if endpoint == "" {
		log.Fatalf("EndpointNamePath pattern can't be empty")
	}

	apiEndpoint := &ApiEndpoint{
		Version:      version,
		EndpointName: endpoint,
	}
	fullPath := apiEndpoint.GetEndpointPath()
	apiEndpoint.Path = fullPath

	return apiEndpoint
}

// GetEndpointPath Returns the string pattern of the EndpointNamePath with the Version
func (a *ApiEndpoint) GetEndpointPath() string {
	if a.Version != "" {
		return fmt.Sprintf("/%s/%s/", a.Version, a.EndpointName)
	}
	return fmt.Sprintf("/%s", a.EndpointName)
}
