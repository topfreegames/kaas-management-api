package kaas

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
)

type ClusterInfrastructure struct {
	Provider string
}

// GetClusterInfrastructure returns a cluster infrastructure resource in a generic format using the ClusterInfrastructure struct
func GetClusterInfrastructure(infrastructureKind string) (*ClusterInfrastructure, error) {
	var infrastructure *ClusterInfrastructure

	switch infrastructureKind {
	case "DockerCluster":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		infrastructure = &ClusterInfrastructure{
			Provider: "docker",
		}
		return infrastructure, nil
	case "KopsAWSCluster":
		infrastructure = &ClusterInfrastructure{
			Provider: "kops",
		}
		return infrastructure, nil
	case "KopsControlPlane":
		infrastructure = &ClusterInfrastructure{
			Provider: "kops",
		}
		return infrastructure, nil
	}
	return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", infrastructureKind))
}
