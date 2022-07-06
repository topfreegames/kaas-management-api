package kops

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	clusterapikopsv1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKopsMachinePool Returns a KopsMachinePool CR from a specific cluster
func GetKopsMachinePool(k *k8s.Kubernetes, clusterName string, infrastructureName string) (*clusterapikopsv1alpha1.KopsMachinePool, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(k8s.KopsMachinePoolSchemaV1alpha1)

	namespace := k8s.GetClusterNamespace(clusterName)
	kopsMachinePoolRaw, err := resource.Namespace(namespace).Get(context.TODO(), infrastructureName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested KopsMachinePool %s was not found in namespace %s!", infrastructureName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting kopsmachinepool from Kubernetes API: %s\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
	}

	var kopsMachinePool clusterapikopsv1alpha1.KopsMachinePool
	kopsMachinePoolRawJson, err := kopsMachinePoolRaw.MarshalJSON()
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidResource, "could not Marshal kopsmachinepool response")
	}

	err = json.Unmarshal(kopsMachinePoolRawJson, &kopsMachinePool)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidResource, "could not Unmarshal kopsmachinepool JSON into clusterAPI")
	}

	return &kopsMachinePool, nil
}
