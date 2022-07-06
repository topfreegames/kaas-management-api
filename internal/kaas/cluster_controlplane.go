package kaas

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
)

type ClusterControlPlane struct {
	Provider string
}

// TODO Change to get a Cluster CR as parameter, validate it and return all desired CP info
// GetControlPlane returns a Control Plane resource in a generic format using the ClusterControlPlane struct
func GetControlPlane(controlPlaneKind string) (*ClusterControlPlane, error) {
	var controlPlane *ClusterControlPlane

	switch controlPlaneKind {
	case "KubeadmControlPlane":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		controlPlane = &ClusterControlPlane{
			Provider: "kubeadm",
		}
		return controlPlane, nil
	case "KopsControlPlane":
		controlPlane = &ClusterControlPlane{
			Provider: "kops",
		}
		return controlPlane, nil
	}

	return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", controlPlaneKind))
}
