package controller

import (
	"github.com/topfreegames/kaas-management-api/internal/kaas"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
)

// ClusterHandler godoc
// @Summary      Get a cluster
// @Description  Get cluster by the full name and show its configuration
// @Tags         Cluster
// @Accept       json
// @Produce      json
// @Param        clusterName   path      string  true  "Cluster Name"
// @Success      200  {object}  v1.Cluster
// @Failure      404  {object}  error.ClientErrorResponse
// @Failure      500  {object}  error.ClientErrorResponse
// @Router       /v1/clusters/{clusterName}/ [get]
// @Security BasicAuth
func (controller ControllerConfig) ClusterHandler(c *gin.Context) {
	clusterName := c.Param(v1.ClusterNameParameter)

	cluster, err := kaas.GetCluster(controller.K8sInstance, clusterName)
	if err != nil {
		log.Printf("[ClusterHandler] Error getting Cluster: %s", err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			} else if clientErr.ErrorMessage == clientError.InvalidConfiguration {
				clientError.ErrorHandler(c, err, clientErr.ErrorDetailedMessage, http.StatusInternalServerError)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	clusterResponse := writeClusterV1Response(cluster)
	c.JSON(http.StatusOK, clusterResponse)
}

// ClusterListHandler godoc
// @Summary      List clusters
// @Description  Return a list of clusters with their information
// @Tags         Cluster
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.ClusterList
// @Failure      404  {object}  error.ClientErrorResponse
// @Failure      500  {object}  error.ClientErrorResponse
// @Router       /v1/clusters/ [get]
// @Security BasicAuth
func (controller ControllerConfig) ClusterListHandler(c *gin.Context) {
	var clusterListResponse v1.ClusterList

	clusterList, err := kaas.ListClusters(controller.K8sInstance)
	if err != nil {
		log.Printf("[ClusterListHandler] Error getting cluster List: %s", err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "No Clusters were found", http.StatusNotFound)
			} else if clientErr.ErrorMessage == clientError.EmptyResponse {
				clientError.ErrorHandler(c, err, "No clusters were found", http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	for _, cluster := range clusterList {
		clusterResponse := writeClusterV1Response(cluster)
		clusterListResponse.Items = append(clusterListResponse.Items, clusterResponse)
	}

	if len(clusterListResponse.Items) == 0 {
		err := clientError.NewClientError(nil, clientError.EmptyResponse, "No Clusters were found")
		clientErr := err.(*clientError.ClientError)
		clientError.ErrorHandler(c, err, clientErr.ErrorDetailedMessage, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, clusterListResponse)
}

// writeClusterV1Response Write the response of the cluster version 1 endpoint
func writeClusterV1Response(cluster *kaas.Cluster) v1.Cluster {
	clusterResponse := v1.Cluster{
		Name:      cluster.Name,
		ApiServer: cluster.ApiEndpoint,
		Metadata: map[string]interface{}{ // TODO, load this mapping from config (kubernetes CR for the management API)
			"clusterGroup": cluster.ClusterGroup,
			"region":       cluster.Region,
			"environment":  cluster.Environment,
			"CIDR":         cluster.CIDR,
		},
		KubeProvider:           cluster.ControlPlane.Provider,
		InfrastructureProvider: cluster.Infrastructure.Provider,
	}
	return clusterResponse
}
