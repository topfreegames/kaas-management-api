package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// GetMachineDeployment Returns a MachineDeployment CR from a specific cluster
func (k Kubernetes) GetMachineDeployment(clusterName string, machineDeploymentName string) (*clusterapiv1beta1.MachineDeployment, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(MachineDeploymentSchemaV1beta1)
	namespace := GetClusterNamespace(clusterName)
	machineDeploymentRaw, err := resource.Namespace(namespace).Get(context.TODO(), machineDeploymentName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested machinedeployment %s was not found for the cluster %s!", machineDeploymentName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting machinedeployment from Kubernetes API: %s\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
	}

	var machineDeployment clusterapiv1beta1.MachineDeployment
	machineDeploymentRawJson, err := machineDeploymentRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinedeployment response: %v", err)
	}

	err = json.Unmarshal(machineDeploymentRawJson, &machineDeployment)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal machinedeployment JSON into clusterAPI list: %v", err)
	}

	err = ValidateMachineTemplateComponents(machineDeployment.Spec.Template)
	if err != nil {
		return nil, clientError.NewClientError(err, clientError.InvalidConfiguration, fmt.Sprintf("MachineDeployment %s doesn't have a valid configuration", machineDeployment.Name))
	}

	return &machineDeployment, nil
}

// ListMachineDeployment Show a list of MachineDeployment CR from a specific cluster
func (k Kubernetes) ListMachineDeployment(clusterName string) (*clusterapiv1beta1.MachineDeploymentList, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(MachineDeploymentSchemaV1beta1)
	namespace := GetClusterNamespace(clusterName)
	machineDeploymentsRaw, err := resource.Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("No machineDeployment was not found for the cluster %s!", clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting machinedeployment list from Kubernetes API: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Kube go-client Error: %v\n", err)
	}

	var machineDeployments clusterapiv1beta1.MachineDeploymentList
	machineDeploymentsRawJson, err := machineDeploymentsRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinedeployment response: %v", err)
	}

	err = json.Unmarshal(machineDeploymentsRawJson, &machineDeployments)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal machinedeployment JSON into clusterAPI list: %v", err)
	}

	if len(machineDeployments.Items) == 0 {
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, fmt.Sprintf("no Machinedeployments were found for the cluster %s!", clusterName))
	}

	return &machineDeployments, nil
}
