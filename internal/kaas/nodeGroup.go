package kaas

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"strings"
)

type NodeGroup struct {
	Name               string
	Cluster            string
	Environment        string
	Region             string
	InfrastructureName string
	InfrastructureKind string
	Replicas           *int32
	Infrastructure     *NodeInfrastructure
}

// GetNodeGroupFullName Returns the real nodeGroup name stored in Kubernetes with the cluster name prefix
func GetNodeGroupFullName(clusterName string, nodeGroupName string) string {
	return fmt.Sprintf("%s-%s", clusterName, nodeGroupName)
}

// GetNodeGroupShortName Returns the nodeGroup name only used by the management API without the cluster name prefix
func GetNodeGroupShortName(clusterName string, nodeGroupFullName string) string {
	return strings.ReplaceAll(nodeGroupFullName, fmt.Sprintf("%s-", clusterName), "")
}

// GetNodeGroup checks which CRD the cluster is using for its node groups (eg machinepool or machinedeployment) and returns a specific node group in the Nodegroup struct format
func GetNodeGroup(k *k8s.Kubernetes, clusterName string, nodeGroupName string) (*NodeGroup, error) {

	nodeGroup := &NodeGroup{
		Name:    nodeGroupName,
		Cluster: clusterName,
	}

	err := nodeGroup.getNodeGroupConfig(k)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroup %s config", nodeGroupName))
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				return nil, clienterr
			} else if clienterr.ErrorMessage == clientError.InvalidConfiguration {
				return nil, clientError.NewClientError(clienterr, clientError.InvalidConfiguration, fmt.Sprintf("NodeGroup %s configuration is invalid", nodeGroupName))
			}
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroup %s config", nodeGroupName))
		}
	}

	infrastructure, err := nodeGroup.getNodeInfrastructure(k)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroup %s infrastructure config", nodeGroupName))
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				return nil, clientError.NewClientError(clienterr, clientError.InvalidResource, fmt.Sprintf("NodeGroup %s is invalid, no infrastructure resource was found or %s.", nodeGroupName, nodeGroup.InfrastructureName))
			} else if clienterr.ErrorMessage == clientError.KindNotFound {
				return nil, clientError.NewClientError(clienterr, clientError.InvalidConfiguration, fmt.Sprintf("NodeGroup %s is invalid, the infrastructure kind %s is not supported.", nodeGroup.InfrastructureKind, nodeGroupName))
			}
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroup %s infrastructure config", nodeGroupName))
		}
	}
	nodeGroup.Infrastructure = infrastructure
	return nodeGroup, nil
}

// getNodeGroupConfig returns the machinePool or machineDeployment configurations used by the nodeGroup.
func (ng *NodeGroup) getNodeGroupConfig(k *k8s.Kubernetes) error {
	// Check if is machinePool
	machinePool, machinePoolErr := k.GetMachinePool(ng.Cluster, GetNodeGroupFullName(ng.Cluster, ng.Name))
	if machinePoolErr != nil {
		clientErr, ok := machinePoolErr.(*clientError.ClientError)
		if !ok {
			return fmt.Errorf("failed getting MachinePool for node group %s in cluster %s: %s", ng.Name, ng.Cluster, machinePoolErr.Error())
		}
		if clientErr.ErrorMessage != clientError.ResourceNotFound {
			return clientError.NewClientError(clientErr, clientError.InvalidConfiguration, fmt.Sprintf("MachinePool %s configuration is invalid", ng.Name))
		}
	} else {
		ng.Cluster = machinePool.Spec.ClusterName
		ng.InfrastructureKind = machinePool.Spec.Template.Spec.InfrastructureRef.Kind
		ng.InfrastructureName = machinePool.Spec.Template.Spec.InfrastructureRef.Name
		ng.Replicas = machinePool.Spec.Replicas
		return nil
	}

	machineDeployment, machineDeploymentErr := k.GetMachineDeployment(ng.Cluster, GetNodeGroupFullName(ng.Cluster, ng.Name))
	if machineDeploymentErr != nil {
		clientErr, ok := machineDeploymentErr.(*clientError.ClientError)
		if !ok {
			return fmt.Errorf("failed getting MachineDeployment for node group %s in cluster %s: %s", ng.Name, ng.Cluster, machinePoolErr.Error())
		}
		if clientErr.ErrorMessage != clientError.ResourceNotFound {
			return clientError.NewClientError(clientErr, clientError.InvalidConfiguration, fmt.Sprintf("MachineDeployment %s configuration is invalid", ng.Name))
		}
	} else {
		ng.InfrastructureKind = machineDeployment.Spec.Template.Spec.InfrastructureRef.Kind
		ng.InfrastructureName = machineDeployment.Spec.Template.Spec.InfrastructureRef.Name
		ng.Replicas = machineDeployment.Spec.Replicas
		return nil
	}

	finalError := fmt.Errorf("Could not get config in neither MachinePool or MachineDeployment: %s, %s", machinePoolErr.Error(), machineDeploymentErr.Error())
	return clientError.NewClientError(finalError, clientError.ResourceNotFound, fmt.Sprintf("Could not find the NodeGroup %s in the cluster %s", ng.Name, ng.Cluster))
}

// ListNodeGroups Returns a list in the Nodegroup struct format
func ListNodeGroups(k *k8s.Kubernetes, clusterName string) ([]*NodeGroup, error) {

	var (
		nodeGroups []*NodeGroup
		hasErrors  bool
	)

	nodeGroupsConfigs, err := GetNodeGroupListConfig(k, clusterName)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroups configurations for cluster %s", clusterName))
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				return nil, clienterr
			} else if clienterr.ErrorMessage == clientError.EmptyResponse {
				return nil, clienterr
			}
			return nil, clientError.NewClientError(clienterr, clientError.UnexpectedError, fmt.Sprintf("Something went wrong while getting NodeGroup configurations for cluster %s", clusterName))
		}
	}

	for _, nodeGroup := range nodeGroupsConfigs {
		infrastructure, err := nodeGroup.getNodeInfrastructure(k)
		if err != nil {
			log.Printf("Error getting NodeInfrastructure for nodegroup %s: %s", nodeGroup.Name, err.Error())
			hasErrors = true
		} else {
			nodeGroup.Infrastructure = infrastructure
			nodeGroups = append(nodeGroups, nodeGroup)
		}
	}

	if len(nodeGroups) < 1 {
		if hasErrors {
			return nil, clientError.NewClientError(nil, clientError.EmptyResponse, fmt.Sprintf("No valid NodeGroups were found for cluster %s, some nodeGroups reported infrastructure resource errors", clusterName))
		}
		return nil, clientError.NewClientError(nil, clientError.EmptyResponse, fmt.Sprintf("No NodeGroups were found for cluster %s", clusterName))
	}

	return nodeGroups, nil
}

// GetNodeGroupListConfig returns the machinePool or machineDeployment configurations used by each nodeGroup.
func GetNodeGroupListConfig(k *k8s.Kubernetes, clusterName string) ([]*NodeGroup, error) {

	var nodeGroups []*NodeGroup
	var validationErr error

	nodePoolErr := map[string]error{
		"machineDeploymentErr": nil,
		"machinePoolErr":       nil,
	}

	// Check if is machinePool
	machinePools, machinePoolErr := k.ListMachinePool(clusterName)
	if machinePoolErr != nil {
		clientErr, ok := machinePoolErr.(*clientError.ClientError)
		if !ok {
			nodePoolErr["machinePoolErr"] = clientError.NewClientError(machinePoolErr, clientError.UnexpectedError, fmt.Sprintf("Error while listing MachinePool for all NodeGroups of the cluster %s", clusterName))
		} else {
			if clientErr.ErrorMessage != clientError.EmptyResponse {
				nodePoolErr["machinePoolErr"] = clientError.NewClientError(clientErr, clientError.UnexpectedError, fmt.Sprintf("Error while listing MachinePool for all NodeGroups of the cluster %s", clusterName))
			}
		}
	} else {
		if len(machinePools.Items) != 0 {
			for _, machinePool := range machinePools.Items {
				validationErr = k8s.ValidateMachineTemplateComponents(machinePool.Spec.Template)
				if validationErr != nil {
					log.Printf("Skipping invalid MachinePool %s: %s", machinePool.Name, validationErr.Error())
					continue
				}
				nodeGroup := &NodeGroup{
					Name:               GetNodeGroupShortName(machinePool.Spec.ClusterName, machinePool.Name),
					Cluster:            machinePool.Spec.ClusterName,
					InfrastructureKind: machinePool.Spec.Template.Spec.InfrastructureRef.Kind,
					InfrastructureName: machinePool.Spec.Template.Spec.InfrastructureRef.Name,
					Replicas:           machinePool.Spec.Replicas,
				}
				nodeGroups = append(nodeGroups, nodeGroup)
			}

			if len(nodeGroups) == 0 {
				return nil, clientError.NewClientError(validationErr, clientError.EmptyResponse, fmt.Sprintf("No valid NodeGroups were found in the cluster %v, some Nodegroups have invalid configuration", clusterName))
			}
			if nodePoolErr["machinePoolErr"] == nil {
				return nodeGroups, nil
			}
		}
	}

	machineDeployments, machineDeploymentErr := k.ListMachineDeployment(clusterName)
	if machineDeploymentErr != nil {
		clientErr, ok := machineDeploymentErr.(*clientError.ClientError)
		if !ok {
			nodePoolErr["machineDeploymentErr"] = clientError.NewClientError(clientErr, clientError.UnexpectedError, fmt.Sprintf("Error while listing MachineDeployment for all NodeGroups of the cluster %s", clusterName))
		} else {
			if clientErr.ErrorMessage != clientError.EmptyResponse {
				nodePoolErr["machineDeploymentErr"] = clientError.NewClientError(machineDeploymentErr, clientErr.ErrorMessage, fmt.Sprintf("Error while listing MachineDeployment for all NodeGroups of the cluster %s", clusterName))
			}
		}
	} else {
		if len(machineDeployments.Items) != 0 {
			for _, machineDeployment := range machineDeployments.Items {
				validationErr = k8s.ValidateMachineTemplateComponents(machineDeployment.Spec.Template)
				if validationErr != nil {
					log.Printf("Skipping invalid MachineDeployment %s: %s", machineDeployment.Name, validationErr.Error())
					continue
				}
				nodeGroup := &NodeGroup{
					Name:               GetNodeGroupShortName(machineDeployment.Spec.ClusterName, machineDeployment.Name),
					Cluster:            machineDeployment.Spec.ClusterName,
					InfrastructureKind: machineDeployment.Spec.Template.Spec.InfrastructureRef.Kind,
					InfrastructureName: machineDeployment.Spec.Template.Spec.InfrastructureRef.Name,
					Replicas:           machineDeployment.Spec.Replicas,
				}
				nodeGroups = append(nodeGroups, nodeGroup)
			}

			if len(nodeGroups) == 0 {
				return nil, clientError.NewClientError(nil, clientError.EmptyResponse, fmt.Sprintf("No valid NodeGroups were found in the cluster %s, some Nodegroups have invalid configuration", clusterName))
			}

			if nodePoolErr["machineDeploymentErr"] == nil {
				return nodeGroups, nil
			}

			return nodeGroups, nil
		}
	}

	if nodePoolErr["machineDeploymentErr"] != nil || nodePoolErr["machinePoolErr"] != nil {
		finalErr := fmt.Errorf(nodePoolErr["machineDeploymentErr"].Error() + " | " + nodePoolErr["machinePoolErr"].Error())
		return nil, clientError.NewClientError(finalErr, clientError.UnexpectedError, fmt.Sprintf("Error while listing infrastructure resources for cluster %s", clusterName))
	}

	return nil, clientError.NewClientError(fmt.Errorf("no nodegroup infrastructure found"), clientError.EmptyResponse, fmt.Sprintf("No NodeGroups were found in the cluster %s", clusterName))
}
