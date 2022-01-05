package controller

import (
	"encoding/json"
	"fmt"
	apiError "github.com/topfreegames/kaas-management-api/api/error"
	nodegroupv1 "github.com/topfreegames/kaas-management-api/api/nodeGroup/v1"
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	clusterv1 "github.com/topfreegames/kaas-management-api/api/cluster/v1"
	"github.com/topfreegames/kaas-management-api/test"
)

func Test_NodeGroupByClusterHandler_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "Success getting nodeGroup in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: nodegroupv1.NodeGroup{
					Name: "nodes",
					Metadata: &nodegroupv1.Metadata{
						Cluster:     "test-cluster.cluster.example.com",
						Replicas:    nil,
						MachineType: "m5.xlarge",
						Zones:       []string{"us-east-1a"},
						Environment: "test",
						Region:      "us-east-1",
						Min:         nil,
						Max:         nil,
					},
					InfrastructureProvider: "kops",
				},
				ExpectedCode: http.StatusOK,
			},
			ExpectedHTTPError: nil,
			Request: &test.HTTPTestRequest{
				Method: "GET",
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups/nodes",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster.cluster.example.com-nodes", "test-cluster.cluster.example.com", "KopsMachinePool", "test-cluster.cluster.example.com-TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestKopsMachinePool("test-cluster.cluster.example.com-TestKopsMachinePool", "test-cluster.cluster.example.com"),
			},
		},
	}

	k := &k8s.Kubernetes{
		K8sAuth: &k8s.Auth{
			DynamicClient: test.NewK8sFakeDynamicClient(),
		},
	}
	controller := ConfigureControllers(k)
	endpoint := test.SetupEndpointRouter(clusterv1.Endpoint)
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf(":%s/nodegroups/:%s", clusterv1.ClusterNameParameter, nodegroupv1.NodeGroupNameParameter), controller.NodeGroupByClusterHandler)

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		request := testCase.GetHTTPRequest()
		expectedResponse, ok := testCase.ExpectedSuccess.(test.HTTPTestExpectedResponse)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *test.HTTPTestExpectedResponse", testCase.Name)
		}

		t.Run(testCase.Name, func(t *testing.T) {
			w := request.RunHTTPTest(endpoint.Router)
			assert.Equal(t, expectedResponse.ExpectedCode, w.Code)
			expected, err := json.Marshal(expectedResponse.ExpectedBody)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}

func Test_NodeGroupByClusterHandler_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Error getting non-existent nodeGroup in clusterV1 endpoint should return not found",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Nodegroup not found",
				ErrorType:    clientError.ResourceNotFound,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups/non-existent",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "Error getting nodeGroup for non-existent cluster in clusterV1 endpoint should return not found",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Cluster not found",
				ErrorType:    clientError.ResourceNotFound,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/non-existent-cluster/nodegroups/nodes",
			},
			K8sTestResources: []runtime.Object{},
		},
		{
			Name:            "Error getting nodeGroup without infrastructure in clusterV1 endpoint should return invalid resource",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Nodegroup configuration is invalid",
				ErrorType:    clientError.InvalidConfiguration,
				HttpCode:     http.StatusInternalServerError,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups/nodes",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster.cluster.example.com-nodes", "test-cluster.cluster.example.com", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "Error getting nodeGroup with invalid infrastructure kind in clusterV1 endpoint should return invalid resource",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Nodegroup configuration is invalid",
				ErrorType:    clientError.InvalidConfiguration,
				HttpCode:     http.StatusInternalServerError,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups/nodes",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster.cluster.example.com-nodes", "test-cluster.cluster.example.com", "invalidKind", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	k := &k8s.Kubernetes{
		K8sAuth: &k8s.Auth{
			DynamicClient: test.NewK8sFakeDynamicClient(),
		},
	}
	controller := ConfigureControllers(k)
	endpoint := test.SetupEndpointRouter(clusterv1.Endpoint)
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf(":%s/nodegroups/:%s", clusterv1.ClusterNameParameter, nodegroupv1.NodeGroupNameParameter), controller.NodeGroupByClusterHandler)

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)
		request := testCase.GetHTTPRequest()

		t.Run(testCase.Name, func(t *testing.T) {
			w := request.RunHTTPTest(endpoint.Router)
			assert.Equal(t, testCase.ExpectedHTTPError.HttpCode, w.Code)
			expected, err := json.Marshal(testCase.ExpectedHTTPError)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}

func Test_NodeGroupListByClusterHandler_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "Success listing one nodegroup in the cluster test-cluster in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: nodegroupv1.NodeGroupList{
					Items: []nodegroupv1.NodeGroup{
						{
							Name: "nodes",
							Metadata: &nodegroupv1.Metadata{
								Cluster:     "test-cluster.cluster.example.com",
								Replicas:    nil,
								MachineType: "m5.xlarge",
								Zones:       []string{"us-east-1a"},
								Environment: "test",
								Region:      "us-east-1",
								Min:         nil,
								Max:         nil,
							},
							InfrastructureProvider: "kops",
						},
					},
				},
				ExpectedCode: http.StatusOK,
			},
			ExpectedHTTPError: nil,
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster.cluster.example.com-nodes", "test-cluster.cluster.example.com", "KopsMachinePool", "test-cluster.cluster.example.com-TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestKopsMachinePool("test-cluster.cluster.example.com-TestKopsMachinePool", "test-cluster.cluster.example.com"),
			},
		},
		{
			Name: "Success listing two nodegroups in the cluster test-cluster2 in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: nodegroupv1.NodeGroupList{
					Items: []nodegroupv1.NodeGroup{
						{
							Name: "nodes2",
							Metadata: &nodegroupv1.Metadata{
								Cluster:     "test-cluster2.cluster.example.com",
								Replicas:    nil,
								MachineType: "m5.xlarge",
								Zones:       []string{"us-east-1a"},
								Environment: "test",
								Region:      "us-east-1",
								Min:         nil,
								Max:         nil,
							},
							InfrastructureProvider: "kops",
						},
						{
							Name: "nodes3",
							Metadata: &nodegroupv1.Metadata{
								Cluster:     "test-cluster2.cluster.example.com",
								Replicas:    nil,
								MachineType: "m5.xlarge",
								Zones:       []string{"us-east-1a"},
								Environment: "test",
								Region:      "us-east-1",
								Min:         nil,
								Max:         nil,
							},
							InfrastructureProvider: "kops",
						},
					},
				},
				ExpectedCode: http.StatusOK,
			},
			ExpectedHTTPError: nil,
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster2.cluster.example.com/nodegroups",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster2.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster2.cluster.example.com-nodes2", "test-cluster2.cluster.example.com", "KopsMachinePool", "test-cluster2.cluster.example.com-TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster2.cluster.example.com-nodes3", "test-cluster2.cluster.example.com", "KopsMachinePool", "test-cluster2.cluster.example.com-TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestKopsMachinePool("test-cluster2.cluster.example.com-TestKopsMachinePool", "test-cluster2.cluster.example.com"),
			},
		},
	}

	k := &k8s.Kubernetes{
		K8sAuth: &k8s.Auth{
			DynamicClient: test.NewK8sFakeDynamicClient(),
		},
	}
	controller := ConfigureControllers(k)
	endpoint := test.SetupEndpointRouter(clusterv1.Endpoint)
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf("/:%s/nodegroups", clusterv1.ClusterNameParameter), controller.NodeGroupListByClusterHandler)

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		request := testCase.GetHTTPRequest()
		expectedResponse, ok := testCase.ExpectedSuccess.(test.HTTPTestExpectedResponse)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *test.HTTPTestExpectedResponse", testCase.Name)
		}

		t.Run(testCase.Name, func(t *testing.T) {
			w := request.RunHTTPTest(endpoint.Router)
			assert.Equal(t, expectedResponse.ExpectedCode, w.Code)
			expected, err := json.Marshal(expectedResponse.ExpectedBody)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}

func Test_NodeGroupListByClusterHandler_ErrorEmptyResponse(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Error getting Nodegroup list for non-existent cluster in clusterV1 endpoint should return not-found",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Cluster not found",
				ErrorType:    clientError.ResourceNotFound,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com/nodegroups",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster2.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster2.cluster.example.com-nodes2", "test-cluster2.cluster.example.com", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("test-cluster2.cluster.example.com-nodes3", "test-cluster2.cluster.example.com", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestKopsMachinePool("test-cluster2.cluster.example.com-TestKopsMachinePool", "test-cluster2.cluster.example.com"),
			},
		},
		{
			Name:            "Error getting Nodegroup list for cluster without nodegroups in clusterV1 endpoint should return empty response",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "No Nodegroups were found for the cluster test-cluster3.cluster.example.com",
				ErrorType:    clientError.EmptyResponse,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster3.cluster.example.com/nodegroups",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster3.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestCluster("test-cluster2.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("test-cluster2.cluster.example.com-nodes2", "test-cluster2.cluster.example.com", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("test-cluster2.cluster.example.com-nodes3", "test-cluster2.cluster.example.com", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestKopsMachinePool("test-cluster2.cluster.example.com-TestKopsMachinePool", "test-cluster2.cluster.example.com"),
			},
		},
	}

	k := &k8s.Kubernetes{
		K8sAuth: &k8s.Auth{
			DynamicClient: test.NewK8sFakeDynamicClient(),
		},
	}

	controller := ConfigureControllers(k)
	endpoint := test.SetupEndpointRouter(clusterv1.Endpoint)
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf("/:%s/nodegroups", clusterv1.ClusterNameParameter), controller.NodeGroupListByClusterHandler)

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)
		request := testCase.GetHTTPRequest()

		t.Run(testCase.Name, func(t *testing.T) {
			w := request.RunHTTPTest(endpoint.Router)
			assert.Equal(t, testCase.ExpectedHTTPError.HttpCode, w.Code)
			expected, err := json.Marshal(testCase.ExpectedHTTPError)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}
