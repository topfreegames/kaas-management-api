package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

// GetMachinePool Returns a MachinePool CR from a specific cluster
func (k Kubernetes) GetMachinePool(clusterName string, machinePoolName string) (*clusterapiexpv1beta1.MachinePool, error) {

	client := k.K8sAuth.DynamicClient

	resource := client.Resource(MachinePoolSchemaV1beta1)
	namespace := GetClusterNamespace(clusterName)
	machinePoolRaw, err := resource.Namespace(namespace).Get(context.TODO(), machinePoolName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested MachinePool %s was not found for the cluster %s!", machinePoolName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting MachinePool from Kubernetes API: %s\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
	}

	var machinePool clusterapiexpv1beta1.MachinePool
	machinePoolRawJson, err := machinePoolRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinePool response: %v", err)
	}

	err = json.Unmarshal(machinePoolRawJson, &machinePool)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal MachinePool JSON into clusterAPI list: %v", err)
	}

	err = ValidateMachineTemplateComponents(machinePool.Spec.Template)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidConfiguration, fmt.Sprintf("MachinePool %s doesn't have a valid configuration", machinePool.Name))
	}

	return &machinePool, nil
}

// ListMachinePool Show a list of MachinePool CR from a specific cluster
func (k Kubernetes) ListMachinePool(clusterName string) (*clusterapiexpv1beta1.MachinePoolList, error) {

	client := k.K8sAuth.DynamicClient

	resource := client.Resource(MachinePoolSchemaV1beta1)
	namespace := GetClusterNamespace(clusterName)
	machinePoolsRaw, err := resource.Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("no MachinePools were found for the cluster %s!", clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting MachinePool list from Kubernetes API: %s\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
	}

	var machinePools clusterapiexpv1beta1.MachinePoolList
	machinePoolsRawJson, err := machinePoolsRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal MachinePool response: %v", err)
	}

	err = json.Unmarshal(machinePoolsRawJson, &machinePools)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal MachinePool JSON into clusterAPI list: %v", err)
	}

	if len(machinePools.Items) == 0 {
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, fmt.Sprintf("no MachinePools were found for the cluster %s!", clusterName))
	}

	return &machinePools, nil
}
