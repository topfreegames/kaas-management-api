package cluster

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	v1 "github.com/topfreegames/kaas-management-api/apis/cluster/v1"
	"github.com/topfreegames/kaas-management-api/test"
	"net/http"
	"testing"
)

func TestClusterHandler(t *testing.T) {
	tests := map[string]test.Case{
		"Success getting test-cluster": {
			ExpectedCode: http.StatusOK,
			ExpectedBody: v1.Cluster{
				Name: "test-cluster.sa-east-1.k8s.tfgco.com",
				Metadata: map[string]interface{}{
					"clusterGroup": "sre-test-clusters",
					"region":       "sa-east-1",
					"environment":  "test",
					"CIDR":         "192.168.0.0/24",
				},
				KubeProvider:           "kops",
				InfrastructureProvider: "aws",
			},
			Request: &test.Request{
				Method: "GET",
				Body:   nil,
				Path:   "/cluster/test-cluster.sa-east-1.k8s.tfgco.com",
			},
		},
		"Fail to get cluster": {
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: v1.Cluster{},
			Request: &test.Request{
				Method: "GET",
				Body:   nil,
				Path:   "/cluster/does-not-exist",
			},
		},
	}

	router := gin.Default()
	router.GET("/cluster/:name", ClusterHandler)
	for testMsg, testCase := range tests {
		t.Run(testMsg, func(t *testing.T) {
			w := testCase.Request.Run(router)
			assert.Equal(t, testCase.ExpectedCode, w.Code)
			expected, err := json.Marshal(testCase.ExpectedBody)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}

