package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

type NodeGroup struct {
	Name               string
	Cluster            string
	Environment        string
	Region             string
	InfrastructureName string
	InfrastructureKind string
	Replicas           *int32
}

// GetMachinePool Returns a Machinepool CR from a specific cluster
func (k Kubernetes) GetMachinePool(clusterName string, nodeGroupName string) (*clusterapiexpv1beta1.MachinePool, error) {

	client := k.K8sAuth.DynamicClient

	resource := client.Resource(machinePoolSchemaV1beta1)
	machinePoolRaw, err := resource.Namespace(clusterName).Get(context.TODO(), nodeGroupName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested machinepool %s was not found for the cluster %s!", nodeGroupName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting machinepool: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
	}

	var machinePool clusterapiexpv1beta1.MachinePool
	machinePoolRawJson, err := machinePoolRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinepool response: %v", err)
	}

	err = json.Unmarshal(machinePoolRawJson, &machinePool)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal machinepool JSON into clusterAPI list: %v", err)
	}

	//err := sanitizeMachinePool(&machinePool)
	//if err != nil {
	//	return nil, clientError.NewClientError(err, clientError.InvalidResource, fmt.Sprintf("The requested machinepool have a invalid spec: %v", err))
	//}

	return &machinePool, nil
}

// ListMachinePool Show a list of Machinepool CR from a specific cluster
func (k Kubernetes) ListMachinePool(clusterName string) (*clusterapiexpv1beta1.MachinePoolList, error) {

	client := k.K8sAuth.DynamicClient

	resource := client.Resource(machinePoolSchemaV1beta1)
	machinePoolsRaw, err := resource.Namespace(clusterName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("no Machinepools were found for the cluster %s!", clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting machinepool list: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
	}

	var machinePools clusterapiexpv1beta1.MachinePoolList
	machinePoolsRawJson, err := machinePoolsRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinepool response: %v", err) //TODO
	}

	err = json.Unmarshal(machinePoolsRawJson, &machinePools)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal machinepool JSON into clusterAPI list: %v", err) // TODO
	}

	if len(machinePools.Items) == 0 {
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, fmt.Sprintf("no Machinepools were found for the cluster %s!", clusterName))
	}

	return &machinePools, nil
}

// GetMachineDeployment Returns a MachineDeployment CR from a specific cluster
func (k Kubernetes) GetMachineDeployment(clusterName string, nodeGroupName string) (*clusterapiv1beta1.MachineDeployment, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(machineDeploymentSchemaV1beta1)

	machineDeploymentRaw, err := resource.Namespace(clusterName).Get(context.TODO(), nodeGroupName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("The requested machinedeployment %s was not found for the cluster %s!", nodeGroupName, clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			return nil, fmt.Errorf("Error getting machinedeployment: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
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

	return &machineDeployment, nil
}

// ListMachineDeployment Show a list of MachineDeployment CR from a specific cluster
func (k Kubernetes) ListMachineDeployment(clusterName string) (*clusterapiv1beta1.MachineDeploymentList, error) {
	client := k.K8sAuth.DynamicClient

	resource := client.Resource(machineDeploymentSchemaV1beta1)

	machineDeploymentsRaw, err := resource.Namespace(clusterName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, clientError.NewClientError(err, clientError.ResourceNotFound, fmt.Sprintf("No machineDeployment was not found for the cluster %s!", clusterName))
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			//TODO
			return nil, fmt.Errorf("Error getting machinedeployment list: %v\n", statusError.ErrStatus.Message)
		}
		return nil, fmt.Errorf("Internal server clientError: %v\n", err)
	}

	var machineDeployments clusterapiv1beta1.MachineDeploymentList
	machineDeploymentsRawJson, err := machineDeploymentsRaw.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not Marshal machinedeployment response: %v", err) //TODO
	}

	err = json.Unmarshal(machineDeploymentsRawJson, &machineDeployments)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal machinedeployment JSON into clusterAPI list: %v", err) // TODO
	}

	if len(machineDeployments.Items) == 0 {
		return nil, clientError.NewClientError(err, clientError.EmptyResponse, fmt.Sprintf("no Machinedeployments were found for the cluster %s!", clusterName))
	}

	return &machineDeployments, nil
}

// GetNodeGroup checks which CRD the cluster is using for its node groups (eg machinepool or machinedeployment) and returns a specific node group in the Nodegroup struct format
func (k Kubernetes) GetNodeGroup(clusterName string, nodeGroupName string) (*NodeGroup, error) {

	var nodeGroup *NodeGroup

	// Check if is machinePool
	machinePool, machinePoolErr := k.GetMachinePool(clusterName, nodeGroupName)
	if machinePoolErr != nil {
		_, ok := machinePoolErr.(*clientError.ClientError)
		if !ok {
			return nil, fmt.Errorf("failed getting machinepool for node group %s in cluster %s: %v", nodeGroupName, clusterName, machinePoolErr)
		}
	} else {
		nodeGroup = &NodeGroup{
			Name:               machinePool.Name,
			Cluster:            machinePool.Spec.ClusterName,
			InfrastructureKind: machinePool.Spec.Template.Spec.InfrastructureRef.Kind,
			InfrastructureName: machinePool.Spec.Template.Spec.InfrastructureRef.Name,
			Replicas:           machinePool.Spec.Replicas,
		}
		return nodeGroup, nil
	}

	machineDeployment, machineDeploymentErr := k.GetMachineDeployment(clusterName, nodeGroupName)
	if machineDeploymentErr != nil {
		_, ok := machineDeploymentErr.(*clientError.ClientError)
		if !ok {
			return nil, fmt.Errorf("failed getting machinedeployment for node group %s in cluster %s: %v", nodeGroupName, clusterName, machinePoolErr)
		}
	} else {
		nodeGroup = &NodeGroup{
			Name:               machineDeployment.Name,
			Cluster:            machineDeployment.Spec.ClusterName,
			InfrastructureKind: machineDeployment.Spec.Template.Spec.InfrastructureRef.Kind,
			InfrastructureName: machineDeployment.Spec.Template.Spec.InfrastructureRef.Name,
			Replicas:           machineDeployment.Spec.Replicas,
		}
		return nodeGroup, nil
	}

	finalError := fmt.Errorf("NodePoolError: %v, %v", machinePoolErr, machineDeploymentErr)
	return nil, clientError.NewClientError(finalError, clientError.ResourceNotFound, fmt.Sprintf("Could not find the NodeGroup %v in the cluster %v", nodeGroupName, clusterName))
}

// GetNodeGroup checks which CRD the cluster is using for its node groups (eg machinepool or machinedeployment) and returns a list in the Nodegroup struct format
func (k Kubernetes) ListNodeGroup(clusterName string) ([]NodeGroup, error) {

	var nodeGroups []NodeGroup

	// Check if is machinePool
	machinePools, machinePoolErr := k.ListMachinePool(clusterName)
	if machinePoolErr != nil {
		_, ok := machinePoolErr.(*clientError.ClientError)
		if !ok {
			return nil, machinePoolErr // Something wrong with the cluster //TODO
		}
	} else {
		for _, machinePool := range machinePools.Items {
			nodeGroup := NodeGroup{
				Name:               machinePool.Name,
				Cluster:            machinePool.Spec.ClusterName,
				InfrastructureKind: machinePool.Spec.Template.Spec.InfrastructureRef.Kind,
				InfrastructureName: machinePool.Spec.Template.Spec.InfrastructureRef.Name,
				Replicas:           machinePool.Spec.Replicas,
			}
			nodeGroups = append(nodeGroups, nodeGroup)
		}
		return nodeGroups, nil
	}

	machineDeployments, machineDeploymentErr := k.ListMachineDeployment(clusterName)
	if machineDeploymentErr != nil {
		_, ok := machineDeploymentErr.(*clientError.ClientError)
		if !ok {
			return nil, machineDeploymentErr // TODO
		}
	} else {
		for _, machineDeployment := range machineDeployments.Items {
			nodeGroup := NodeGroup{
				Name:               machineDeployment.Name,
				Cluster:            machineDeployment.Spec.ClusterName,
				InfrastructureKind: machineDeployment.Spec.Template.Spec.InfrastructureRef.Kind,
				InfrastructureName: machineDeployment.Spec.Template.Spec.InfrastructureRef.Name,
				Replicas:           machineDeployment.Spec.Replicas,
			}
			nodeGroups = append(nodeGroups, nodeGroup)
		}
		return nodeGroups, nil
	}

	finalError := fmt.Errorf("NodePoolError: %v, %v", machinePoolErr, machineDeploymentErr)
	return nil, clientError.NewClientError(finalError, clientError.EmptyResponse, fmt.Sprintf("No NodeGroups were found in the cluster %v", clusterName))
}
