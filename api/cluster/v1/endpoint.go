package v1

import "github.com/topfreegames/kaas-management-api/api"

var Endpoint = api.NewApiEndpoint("v1", "clusters")

// Parameters
const (
	ClusterNameParameter = "clusterName"
)