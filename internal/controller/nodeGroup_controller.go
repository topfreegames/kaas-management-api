package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	nodegroupv1 "github.com/topfreegames/kaas-management-api/api/nodeGroup/v1"
	"github.com/topfreegames/kaas-management-api/internal/kaas"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"net/http"
)

// NodeGroupByClusterHandler godoc
// @Summary      Get a specific node group from a cluster
// @Description  Shows the information about a node group of a cluster
// @Tags         Cluster
// @Accept       json
// @Produce      json
// @Param        clusterName   path      string  true  "Cluster Name"
// @Param        nodeGroupName   path      string  true  "Node Group Name"
// @Success      200  {object}  nodegroupv1.NodeGroup
// @Failure      404  {object}  error.ClientErrorResponse
// @Failure      500  {object}  error.ClientErrorResponse
// @Router       /v1/clusters/{clusterName}/nodegroup/{nodeGroupName}/ [get]
// @Security BasicAuth
func (controller ControllerConfig) NodeGroupByClusterHandler(c *gin.Context) {
	clusterName := c.Param(clusterv1.ClusterNameParameter)
	nodeGroupName := c.Param(nodegroupv1.NodeGroupNameParameter)

	cluster, err := kaas.GetCluster(controller.K8sInstance, clusterName)
	if err != nil {
		log.Printf("[NodeGroupByClusterHandler] Error getting clusterAPI CR: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	nodeGroup, err := kaas.GetNodeGroup(controller.K8sInstance, clusterName, nodeGroupName)
	if err != nil {
		log.Printf("[NodeGroupByClusterHandler] Error getting NodeGroup: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Nodegroup not found", http.StatusNotFound)
			} else if clienterr.ErrorMessage == clientError.InvalidResource {
				clientError.ErrorHandler(c, err, "Nodegroup resource is invalid", http.StatusInternalServerError)
			} else if clienterr.ErrorMessage == clientError.InvalidConfiguration {
				clientError.ErrorHandler(c, err, "Nodegroup configuration is invalid", http.StatusInternalServerError)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	nodeGroupV1 := writeNodeGroupV1Response(cluster, nodeGroup)
	c.JSON(http.StatusOK, nodeGroupV1)
}

// NodeGroupListByClusterHandler godoc
// @Summary      List node groups from a cluster
// @Description  List all node groups of a specific cluster with each Node Group information
// @Tags         Cluster
// @Accept       json
// @Produce      json
// @Param        clusterName   path      string  true  "Cluster Name"
// @Success      200  {object}  nodegroupv1.NodeGroupList
// @Failure      500  {object}  error.ClientErrorResponse
// @Router       /v1/clusters/{clusterName}/nodegroups/ [get]
// @Security BasicAuth
func (controller ControllerConfig) NodeGroupListByClusterHandler(c *gin.Context) {
	clusterName := c.Param(clusterv1.ClusterNameParameter)

	var nodegroupV1List nodegroupv1.NodeGroupList

	cluster, err := kaas.GetCluster(controller.K8sInstance, clusterName)
	if err != nil {
		log.Printf("[NodeGroupByClusterHandler] Error getting clusterAPI CR: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	nodeGroups, err := kaas.ListNodeGroups(controller.K8sInstance, clusterName)
	if err != nil {
		log.Printf("[NodeGroupListByClusterHandler] Error Listing Nodegroup: %s", err.Error())
		clienterr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clienterr.ErrorMessage == clientError.EmptyResponse {
				clientError.ErrorHandler(c, err, fmt.Sprintf("No Nodegroups were found for the cluster %s", clusterName), http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	for _, nodeGroup := range nodeGroups {
		nodeGroupV1 := writeNodeGroupV1Response(cluster, nodeGroup)
		nodegroupV1List.Items = append(nodegroupV1List.Items, nodeGroupV1)
	}

	if len(nodegroupV1List.Items) == 0 {
		err := clientError.NewClientError(nil, clientError.EmptyResponse, fmt.Sprintf("No Nodegroups were found for the cluster %s", clusterName))
		clientErr := err.(*clientError.ClientError)
		clientError.ErrorHandler(c, err, clientErr.ErrorDetailedMessage, http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, nodegroupV1List)
}

// writeNodeGroupV1Response Write the response of the nodeGroup version 1 endpoint
func writeNodeGroupV1Response(cluster *kaas.Cluster, nodeGroup *kaas.NodeGroup) nodegroupv1.NodeGroup {
	metadata := &nodegroupv1.Metadata{
		Cluster:     nodeGroup.Cluster,
		Replicas:    nodeGroup.Replicas,
		MachineType: nodeGroup.Infrastructure.MachineType,
		Zones:       nodeGroup.Infrastructure.Az,
		Environment: cluster.Environment,
		Region:      cluster.Region,
		Min:         nodeGroup.Infrastructure.Min,
		Max:         nodeGroup.Infrastructure.Max,
	}
	nodeGroupV1 := nodegroupv1.NodeGroup{
		Name:                   nodeGroup.Name,
		Metadata:               metadata,
		InfrastructureProvider: nodeGroup.Infrastructure.Provider,
	}
	return nodeGroupV1
}
