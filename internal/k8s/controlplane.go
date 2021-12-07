package k8s

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
)

type ControlPlane struct {
	Provider string
}

// GetControlPlane returns a Control Plane resource in a generic format using the ControlPlane struct
func (k Kubernetes) GetControlPlane(controlPlaneKind string) (*ControlPlane, error) {
	var controlPlane *ControlPlane

	switch controlPlaneKind {
	case "KubeadmControlPlane":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		controlPlane = &ControlPlane{
			Provider: "kubeadm",
		}
		return controlPlane, nil
	case "KopsControlPlane":
		controlPlane = &ControlPlane{
			Provider: "kops",
		}
		return controlPlane, nil
	}

	return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", controlPlaneKind))
}
