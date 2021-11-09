package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/topfreegames/kaas-management-api/apis/cluster/v1"
)
// ClusterHandler - returns a cluster status
func ClusterHandler(c *gin.Context) {
	clusterName := c.Param("name")

	if clusterName == "test-cluster.sa-east-1.k8s.tfgco.com" {
		cluster := v1.Cluster{
			Name: "test-cluster.sa-east-1.k8s.tfgco.com",
			Metadata: map[string]interface{}{
				"clusterGroup": "sre-test-clusters",
				"region":       "sa-east-1",
				"environment":  "test",
				"CIDR":         "192.168.0.0/24",
			},
			KubeProvider:           "kops",
			InfrastructureProvider: "aws",
		}
		c.JSON(http.StatusOK, cluster)
	} else {
		c.JSON(http.StatusNotFound, v1.Cluster{})
	}

}

// ClusterListHandler - returns a list of clusters status
func ClusterListHandler(c *gin.Context) {
	var clusterList v1.ClusterList

	cluster1 := v1.Cluster{
		Name: "test-cluster.eu-central-1.k8s.tfgco.com",
		Metadata: map[string]interface{}{
			"clusterGroup": "sre-test-clusters",
			"region":       "eu-central-1",
			"environment":  "test",
			"CIDR":         "192.168.1.0/24",
		},
		KubeProvider:           "kops",
		InfrastructureProvider: "aws",
	}

	cluster2 := v1.Cluster{
		Name: "test-cluster.us-east-1.k8s.tfgco.com",
		Metadata: map[string]interface{}{
			"clusterGroup": "sre-test-clusters",
			"region":       "us-east-1",
			"environment":  "test",
			"CIDR":         "192.168.2.0/24",
		},
		KubeProvider:           "kops",
		InfrastructureProvider: "aws",
	}

	clusterList.Items = append(clusterList.Items, cluster1, cluster2)

	c.JSON(http.StatusOK, clusterList)
}