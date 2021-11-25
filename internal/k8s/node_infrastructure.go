package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	clusterapikopsv1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeInfrastructure struct {
	Provider    string
	Az          []string
	MachineType string
	Min         *int32
	Max         *int32
	Spec        interface{}
}

func (k Kubernetes) GetInfrastructure(clusterName, infrastructureKind string, infrastructureName string) (*Infrastructure, error) {
	var infrastructure *Infrastructure

	switch infrastructureKind {
	case "DockerCluster":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		infrastructure = &Infrastructure{
			Provider: "Docker",
			Az: []string{
				"local",
			},
			MachineType: "container",
			Spec:        nil,
		}

	case "DockerMachineTemplate":
		// DockerMachine api is a test resource for cluster-api, it api code breaks often so there's no reason to really use it.
		// TODO: Fork the official repo, fix the go.mod and implement to be used in our tests
		infrastructure = &Infrastructure{
			Provider: "Docker",
			Az: []string{
				"local",
			},
			MachineType: "container",
			Spec:        nil,
		}

	case "KopsMachinePool":
		kops, err := k.GetKopsMachinePool(clusterName, infrastructureName)
		if err != nil {
			clientErr, ok := err.(*clientError.ClientError)
			if !ok {
				return nil, fmt.Errorf("an error has ocurred while feching kopsmachinepool infrastructure: %v", err)
			}
			return nil, clientError.NewClientError(err, clientErr.ErrorMessage, fmt.Sprintf("Could not retrieve the infrastructure: %v", clientErr.ErrorDetailedMessage))
		}
		infrastructure = &Infrastructure{
			Provider:    "kops",
			Az:          kops.Spec.KopsInstanceGroupSpec.Subnets,
			MachineType: kops.Spec.KopsInstanceGroupSpec.MachineType,
			Min:         kops.Spec.KopsInstanceGroupSpec.MinSize,
			Max:         kops.Spec.KopsInstanceGroupSpec.MaxSize,
			Spec:        kops.Spec,
		}
		return infrastructure, nil
	default:
		return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", infrastructureKind))
	}

	return nil, clientError.NewClientError(nil, clientError.KindNotFound, fmt.Sprintf("The Kind %s could not be found", infrastructureKind))
}

func (k Kubernetes) GetKopsMachinePool(clusterName string, infrastructureName string) (*clusterapikopsv1alpha1.KopsMachinePool, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(kopsMachinePoolSchemaV1alpha1)
	kopsMachinePoolRaw, err := resource.Namespace(clusterName).Get(context.TODO(), infrastructureName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested KopsMachinePool %s was not found in namespace %s!", infrastructureName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting kopsmachinepool: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
	}

	var kopsMachinePool clusterapikopsv1alpha1.KopsMachinePool
	kopsMachinePoolRawJson, err := kopsMachinePoolRaw.MarshalJSON()
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidResource, "could not Marshal kopsmachinepool response")
	}

	err = json.Unmarshal(kopsMachinePoolRawJson, &kopsMachinePool)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidResource, "could not Unmarshal kopsmachinepool JSON into clusterAPI list")
	}

	return &kopsMachinePool, nil
}
