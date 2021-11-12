package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/util"
)

func (controller ControllerConfig) NodeGroupByClusterHandler(c *gin.Context) {
	clusterName := c.Param("clusterName")
	nodeGroupName := c.Param("nodeGroupName")

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

	nodeGrouApiCR, err := controller.K8sInstance.GetNodeGroup(clusterApiCR.Name, nodeGroupName)
}

func (controller ControllerConfig) NodeGroupListByClusterHandler(c *gin.Context) {

}
