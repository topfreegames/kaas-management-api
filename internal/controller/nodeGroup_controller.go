package controller

import (
	"fmt"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	nodegroupv1 "github.com/topfreegames/kaas-management-api/api/nodeGroup/v1"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"net/http"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/gin-gonic/gin"
)

// NodeGroupByClusterHandler Shows the information about a node group of a cluster
func (controller ControllerConfig) NodeGroupByClusterHandler(c *gin.Context) {
	clusterName := c.Param(clusterv1.ClusterNameParameter)
	nodeGroupName := c.Param(nodegroupv1.NodeGroupNameParameter)

	cluster, err := controller.K8sInstance.GetCluster(clusterName)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			log.Printf("Error getting clusterAPI CR: %v", err)
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			log.Printf("Error getting clusterAPI CR: %s: %v", clienterr.ErrorDetailedMessage, clienterr.ErrorCause)
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			}
		}
		log.Printf("[NodeGroupByClusterHandler] %v", err)
		return
	}

	nodeGroup, err := controller.K8sInstance.GetNodeGroup(clusterName, nodeGroupName)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			log.Printf("Error getting NodeGroup: %v", err)
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			log.Printf("Error getting NodeGroup: %s: %v", clienterr.ErrorDetailedMessage, clienterr.ErrorCause)
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Nodegroup not found", http.StatusNotFound)
			}
			if clienterr.ErrorMessage == clientError.InvalidResource {
				clientError.ErrorHandler(c, err, fmt.Sprintf("Nodegroup is invalid: %v", clienterr.ErrorCause.Error()), http.StatusInternalServerError)
			}
		}
		log.Printf("[NodeGroupByClusterHandler] %v", err)
		return
	}
	infra, err := controller.K8sInstance.GetNodeInfrastructure(clusterName, nodeGroup.InfrastructureKind, nodeGroup.InfrastructureName)
	if err != nil {
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			log.Printf("Error getting NodeInfrastructure: %v", err)
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			log.Printf("Error getting NodeGroup: %s: %v", clienterr.ErrorDetailedMessage, clienterr.ErrorCause)
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Nodegroup infrastructure resource is invalid", http.StatusInternalServerError)
			}
			if clienterr.ErrorMessage == clientError.InvalidResource {
				clientError.ErrorHandler(c, err, fmt.Sprintf("Nodegroup infrastructure resource is invalid: %v", clienterr.ErrorCause.Error()), http.StatusInternalServerError)
			}
		}
		log.Printf("[NodeGroupByClusterHandler] %v", err)
		return
	}
	nodeGroup.Infrastructure = infra

	nodeGroupV1 := writeNodeGroupV1(cluster, nodeGroup, infra)
	c.JSON(http.StatusOK, nodeGroupV1)
}

// NodeGroupListByClusterHandler List all node groups of a specific cluster with each Node Group information
func (controller ControllerConfig) NodeGroupListByClusterHandler(c *gin.Context) {
	clusterName := c.Param(clusterv1.ClusterNameParameter)

	var nodegroupV1List nodegroupv1.NodeGroupList

	cluster, err := controller.K8sInstance.GetCluster(clusterName)
	if err != nil {
		log.Printf("[NodeGroupListByClusterHandler] Error getting clusterAPI CR: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			}
		}
		return
	}

	nodeGroups, err := controller.K8sInstance.ListNodeGroup(clusterName)
	if err != nil {
		log.Printf("[NodeGroupListByClusterHandler] Error Listing Nodegroup: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.EmptyResponse {
				clientError.ErrorHandler(c, err, fmt.Sprintf("No Nodegroups were found for the cluster %s", clusterName), http.StatusNotFound)
			}
		}
		return
	}

	for _, nodeGroup := range nodeGroups {
		infra, err := controller.K8sInstance.GetNodeInfrastructure(clusterName, nodeGroup.InfrastructureKind, nodeGroup.InfrastructureName)
		if err != nil {
			log.Printf("[NodeGroupHandler] Error getting NodeInfrastructure for nodegroup %s: %s", nodeGroup.Name, err.Error())
		} else {
			nodeGroup.Infrastructure = infra
			nodeGroupV1 := writeNodeGroupV1(cluster, nodeGroup, infra)
			nodegroupV1List.Items = append(nodegroupV1List.Items, nodeGroupV1)
		}
	}

	if len(nodegroupV1List.Items) == 0 {
		err := clientError.NewClientError(nil, clientError.EmptyResponse, fmt.Sprintf("[NodeGroupHandler] No Nodegroups were found for the cluster %s", clusterName))
		clientErr := err.(*clientError.ClientError)
		clientError.ErrorHandler(c, err, clientErr.ErrorDetailedMessage, http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, nodegroupV1List)

}

// TODO better function name
// writeNodeGroupV1 Write the response of the nodeGroup version 1 endpoint
func writeNodeGroupV1(cluster clusterapiv1beta1.Cluster, nodeGroup *k8s.NodeGroup, nodeGroupInfrastructure *k8s.NodeInfrastructure) nodegroupv1.NodeGroup {

	metadata := &nodegroupv1.Metadata{
		Cluster:     cluster.Name,
		Replicas:    nodeGroup.Replicas,
		MachineType: nodeGroupInfrastructure.MachineType,
		Zones:       nodeGroupInfrastructure.Az,
		Environment: cluster.Labels["environment"],
		Region:      cluster.Labels["region"],
		Min:         nodeGroupInfrastructure.Min,
		Max:         nodeGroupInfrastructure.Max,
	}

	nodeGroupV1 := nodegroupv1.NodeGroup{
		Name:                   nodeGroup.Name,
		Metadata:               metadata,
		KubeProvider:           "",
		InfrastructureProvider: nodeGroupInfrastructure.Provider,
	}

	return nodeGroupV1
}