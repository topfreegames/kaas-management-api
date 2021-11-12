package controller

import (
	"log"
	"net/http"

	"github.com/topfreegames/kaas-management-api/util"

	"sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/gin-gonic/gin"
	v1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
)

// ClusterHandler - returns a cluster status
func (controller ControllerConfig) ClusterHandler(c *gin.Context) {
	clusterName := c.Param("clusterName")

	// TODO this shouldn't receive a namespace
	clusterApiCR, err := controller.K8sInstance.GetCluster(clusterName, "default") // TODO remove hardcode default namespace
	if err != nil {
		log.Printf("Error getting clusterAPI CR: %v", err)
		_, ok := err.(*util.ClientError)
		if !ok {
			util.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
		} else {
			util.ErrorHandler(c, err, "Cluster not found", http.StatusNotFound)
		}
		log.Printf("[ClusterHandler] %v", err)
		return
	}

	cluster := writeClusterV1(clusterApiCR)
	c.JSON(http.StatusOK, cluster)
}

// ClusterListHandler - returns a list of clusters status
func (controller ControllerConfig) ClusterListHandler(c *gin.Context) {
	var clusterList v1.ClusterList

	clusterApiListCR, err := controller.K8sInstance.ListClusters("")
	if err != nil {
		log.Printf("Error getting clusterAPI CR: %v", err)
		_, ok := err.(util.ClientError)
		if !ok {
			log.Printf("[ClusterListHandler] %v", err)
			util.ErrorHandler(c, err, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	for _, clusterApiCR := range clusterApiListCR.Items {
		cluster := writeClusterV1(clusterApiCR)
		clusterList.Items = append(clusterList.Items, cluster)
	}

	// TODO check if we should return StatusOK even with an empty list
	c.JSON(http.StatusOK, clusterList)
}

// TODO better function name
func writeClusterV1(clusterApiCR v1beta1.Cluster) v1.Cluster {
	cluster := v1.Cluster{
		Name: clusterApiCR.Name,
		// TODO, load this mapping from config (kubernetes CR for the management API)
		Metadata: map[string]interface{}{
			"clusterGroup": clusterApiCR.Labels["clusterGroup"],
			"region":       clusterApiCR.Labels["region"],
			"environment":  clusterApiCR.Labels["environment"],
			"CIDR":         clusterApiCR.Spec.ClusterNetwork.Services.CIDRBlocks,
		},
		KubeProvider:           getKubeProvider(clusterApiCR),
		InfrastructureProvider: getInfrastructureProvider(clusterApiCR),
	}
	return cluster
}

// TODO some enum/dict with supported providers and the kind names
func getKubeProvider(cluster v1beta1.Cluster) string {
	controlplaneProvider := cluster.Spec.ControlPlaneRef.Kind

	if controlplaneProvider == "KubeadmControlPlane" {
		return "KubeAdm"
	}

	return "Undefined"
}

// TODO some enum/dict with supported providers and the kind names
func getInfrastructureProvider(cluster v1beta1.Cluster) string {
	controlplaneProvider := cluster.Spec.InfrastructureRef.Kind

	if controlplaneProvider == "DockerCluster" {
		return "Docker"
	}

	return "Undefined"
}
