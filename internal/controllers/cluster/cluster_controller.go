package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/topfreegames/kaas-management-api/apis/cluster/v1"
)

// PingHandler returns pong
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

// ClusterHandler - returns a cluster status
func ClusterHandler(c *gin.Context) {
	cluster := v1.Cluster{
		Name: "test-cluster.sa-east-1.k8s.tfgco.com",
		Metadata: map[string]interface{}{
			"clusterGroup": "sre-test-clusters",
			"region":       "us-east-1",
			"environment":  "test",
			"CIDR":         "192.168.0.0/24",
		},
		KubeProvider:           "kops",
		InfrastructureProvider: "aws",
	}
	c.JSON(http.StatusOK, cluster)
}
