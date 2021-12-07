package controller

import (
	"fmt"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"log"
	"net/http"

	"sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/gin-gonic/gin"
	v1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
)

// ClusterHandler - returns a cluster information
func (controller ControllerConfig) ClusterHandler(c *gin.Context) {
	clusterName := c.Param(v1.ClusterNameParameter)

	clusterApiCR, err := controller.K8sInstance.GetCluster(clusterName)
	if err != nil {
		log.Printf("[ClusterHandler] Error getting clusterAPI CR: %s", err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	controlPlane, err := controller.K8sInstance.GetControlPlane(clusterApiCR.Spec.ControlPlaneRef.Kind)
	if err != nil {
		log.Printf("[ClusterHandler] Error getting cluster controlplane for cluster %s: %s", clusterApiCR.Name, err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.KindNotFound {
				newErr := clientError.NewClientError(err, clientError.InvalidConfiguration, clientErr.ErrorDetailedMessage)
				clientError.ErrorHandler(c, newErr, "Cluster configuration is invalid", http.StatusInternalServerError)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}
	infrastructure, err := controller.K8sInstance.GetClusterInfrastructure(clusterApiCR.Spec.InfrastructureRef.Kind)
	if err != nil {
		log.Printf("Error getting cluster infrastructure for cluster %s: %s", clusterApiCR.Name, err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.KindNotFound {
				newErr := clientError.NewClientError(err, clientError.InvalidConfiguration, clientErr.ErrorDetailedMessage)
				clientError.ErrorHandler(c, newErr, "Cluster configuration is invalid", http.StatusInternalServerError)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	cluster := writeClusterV1(clusterApiCR, controlPlane, infrastructure)
	c.JSON(http.StatusOK, cluster)
}

// ClusterListHandler - returns a list of clusters with their information
func (controller ControllerConfig) ClusterListHandler(c *gin.Context) {
	var clusterList v1.ClusterList

	clusterApiListCR, err := controller.K8sInstance.ListClusters()
	if err != nil {
		log.Printf("[ClusterListHandler] Error getting clusterAPI CR: %s", err.Error())
		clientErr, ok := err.(*clientError.ClientError)
		if !ok {
			clientError.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			if clientErr.ErrorMessage == clientError.ResourceNotFound {
				clientError.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
			} else {
				clientError.ErrorHandler(c, err, "Unhandled Error", http.StatusInternalServerError)
			}
		}
		return
	}

	for _, clusterApiCR := range clusterApiListCR.Items {
		controlPlane, err := controller.K8sInstance.GetControlPlane(clusterApiCR.Spec.ControlPlaneRef.Kind)
		if err != nil {
			log.Printf("Error getting cluster controlplane for cluster %s: %v", clusterApiCR.Name, err.Error())
			continue
		}

		infrastructure, err := controller.K8sInstance.GetClusterInfrastructure(clusterApiCR.Spec.InfrastructureRef.Kind)
		if err != nil {
			log.Printf("Error getting cluster infrastructure for %s: %v", clusterApiCR.Name, err.Error())
			continue
		}
		cluster := writeClusterV1(&clusterApiCR, controlPlane, infrastructure)
		clusterList.Items = append(clusterList.Items, cluster)
	}

	if len(clusterList.Items) == 0 {
		err := clientError.NewClientError(nil, clientError.EmptyResponse, "No Clusters were found")
		clientErr := err.(*clientError.ClientError)
		clientError.ErrorHandler(c, err, clientErr.ErrorDetailedMessage, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, clusterList)
}

// TODO better function name
// writeClusterV1 Write the response of the cluster version 1 endpoint
func writeClusterV1(clusterApiCR *v1beta1.Cluster, controlPlane *k8s.ControlPlane, infrastructure *k8s.ClusterInfrastructure) v1.Cluster {
	apiEndpoint := fmt.Sprintf("https://%s:%d", clusterApiCR.Spec.ControlPlaneEndpoint.Host, clusterApiCR.Spec.ControlPlaneEndpoint.Port)

	cluster := v1.Cluster{
		Name:      clusterApiCR.Name,
		ApiServer: apiEndpoint,
		// TODO, load this mapping from config (kubernetes CR for the management API)
		Metadata: map[string]interface{}{
			"clusterGroup": clusterApiCR.Labels["clusterGroup"],
			"region":       clusterApiCR.Labels["region"],
			"environment":  clusterApiCR.Labels["environment"],
			"CIDR":         clusterApiCR.Spec.ClusterNetwork.Services.CIDRBlocks,
		},
		KubeProvider:           controlPlane.Provider,
		InfrastructureProvider: infrastructure.Provider,
	}
	return cluster
}
