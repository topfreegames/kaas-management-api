package kaas

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/internal/k8s/providers/kops"
	"github.com/topfreegames/kaas-management-api/util/clientError"
)

type NodeInfrastructure struct {
	Name        string
	Cluster     string
	Provider    string
	Az          []string
	MachineType string
	Min         *int32
	Max         *int32
	Spec        interface{}
}

// GetNodeInfrastructure returns a nodegroup infrastructure resource in a generic format using the NodeInfrastructure struct
func (ng *NodeGroup) getNodeInfrastructure(k *k8s.Kubernetes) (*NodeInfrastructure, error) {
	var infrastructure *NodeInfrastructure

	switch ng.InfrastructureKind {
	case "DockerMachineTemplate":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it other than development.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		infrastructure = &NodeInfrastructure{
			Name:     "docker",
			Provider: "docker",
			Cluster:  "docker-cluster",
			Az: []string{
				"local",
			},
			MachineType: "container",
			Spec:        nil,
		}
		return infrastructure, nil

	case "KopsMachinePool":
		kops, err := kops.GetKopsMachinePool(k, ng.Cluster, ng.InfrastructureName)
		if err != nil {
			clientErr, ok := err.(*clientError.ClientError)
			if !ok {
				return nil, fmt.Errorf("an error has ocurred while feching kopsmachinepool infrastructure: %s", err.Error())
			}
			return nil, clientError.NewClientError(clientErr, clientErr.ErrorMessage, "Could not retrieve the infrastructure")
		}
		infrastructure = &NodeInfrastructure{
			Name:        kops.Name,
			Provider:    "kops",
			Cluster:     kops.ClusterName,
			Az:          kops.Spec.KopsInstanceGroupSpec.Subnets,
			MachineType: kops.Spec.KopsInstanceGroupSpec.MachineType,
			Min:         kops.Spec.KopsInstanceGroupSpec.MinSize,
			Max:         kops.Spec.KopsInstanceGroupSpec.MaxSize,
			Spec:        kops.Spec,
		}
		return infrastructure, nil
	}

	return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", ng.InfrastructureKind))
}
