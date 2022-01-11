package controller

import (
	"encoding/json"
	"fmt"
	apiError "github.com/topfreegames/kaas-management-api/api/error"
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

func Test_ClusterHandler_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "Success getting test-cluster in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: clusterv1.Cluster{
					Name:      "test-cluster.cluster.example.com",
					ApiServer: "https://api.test-cluster.cluster.example.com.cluster.example.com:443",
					Metadata: map[string]interface{}{
						"clusterGroup": "test-clusters",
						"region":       "us-east-1",
						"environment":  "test",
						"CIDR":         []string{"192.168.0.0/24"},
					},
					KubeProvider:           "kops",
					InfrastructureProvider: "kops",
				},
				ExpectedCode: http.StatusOK,
			},
			ExpectedHTTPError: nil,
			Request: &test.HTTPTestRequest{
				Method: "GET",
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
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
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf(":%s", clusterv1.ClusterNameParameter), controller.ClusterHandler)

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

func Test_ClusterHandler_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Error getting test-cluster in clusterV1 endpoint should return not found",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Cluster not found",
				ErrorType:    clientError.ResourceNotFound,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com",
			},
			K8sTestResources: []runtime.Object{},
		},
		{
			Name:            "Error getting test-cluster in clusterV1 endpoint should return invalid for non-existent infrastructure kind",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Cluster configuration is invalid",
				ErrorType:    clientError.InvalidConfiguration,
				HttpCode:     http.StatusInternalServerError,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "Error getting test-cluster in clusterV1 endpoint should return invalid for non-existent controlplane kind",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "Cluster configuration is invalid",
				ErrorType:    clientError.InvalidConfiguration,
				HttpCode:     http.StatusInternalServerError,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/test-cluster.cluster.example.com",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "non-existent", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "kops-cluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
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
	endpoint.CreateRoute(http.MethodGet, fmt.Sprintf(":%s", clusterv1.ClusterNameParameter), controller.ClusterHandler)

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

func Test_ClusterListHandler_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "Success getting one cluster in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: clusterv1.ClusterList{
					Items: []clusterv1.Cluster{
						{
							Name:      "test-cluster.cluster.example.com",
							ApiServer: "https://api.test-cluster.cluster.example.com.cluster.example.com:443",
							Metadata: map[string]interface{}{
								"clusterGroup": "test-clusters",
								"region":       "us-east-1",
								"environment":  "test",
								"CIDR":         []string{"192.168.0.0/24"},
							},
							KubeProvider:           "kops",
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
				Path:   clusterv1.Endpoint.EndpointPath + "/",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name: "Success getting two clusters in clusterV1 endpoint",
			ExpectedSuccess: test.HTTPTestExpectedResponse{
				ExpectedBody: clusterv1.ClusterList{
					Items: []clusterv1.Cluster{
						{
							Name:      "test-cluster.cluster.example.com",
							ApiServer: "https://api.test-cluster.cluster.example.com.cluster.example.com:443",
							Metadata: map[string]interface{}{
								"clusterGroup": "test-clusters",
								"region":       "us-east-1",
								"environment":  "test",
								"CIDR":         []string{"192.168.0.0/24"},
							},
							KubeProvider:           "kops",
							InfrastructureProvider: "kops",
						},
						{
							Name:      "test-cluster2.cluster.example.com",
							ApiServer: "https://api.test-cluster2.cluster.example.com.cluster.example.com:443",
							Metadata: map[string]interface{}{
								"clusterGroup": "test-clusters",
								"region":       "us-east-1",
								"environment":  "test",
								"CIDR":         []string{"192.168.0.0/24"},
							},
							KubeProvider:           "kops",
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
				Path:   clusterv1.Endpoint.EndpointPath + "/",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestCluster("test-cluster2.cluster.example.com", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
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
	endpoint.CreateRoute(http.MethodGet, "/", controller.ClusterListHandler)

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

func Test_ClusterListHandler_ErrorEmptyResponse(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Error getting cluster list in clusterV1 endpoint should return empty response for invalid cluster in list",
			ExpectedSuccess: nil,
			ExpectedHTTPError: &apiError.ClientErrorResponse{
				ErrorMessage: "No clusters were found",
				ErrorType:    clientError.EmptyResponse,
				HttpCode:     http.StatusNotFound,
			},
			Request: &test.HTTPTestRequest{
				Method: http.MethodGet,
				Body:   nil,
				Path:   clusterv1.Endpoint.EndpointPath + "/",
			},
			K8sTestResources: []runtime.Object{
				// Invalid cluster without controPlane and infrastructure
				test.NewTestCluster("test-cluster.cluster.example.com", "testcluster-kops-cp", "", "controlplane.cluster.x-k8s.io/v1alpha1", "", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
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
	endpoint.CreateRoute(http.MethodGet, "/", controller.ClusterListHandler)

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
