package kaas

import (
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"reflect"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"testing"
)

func Test_GetCluster_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetCluster should return Success for one cluster",
			ExpectedSuccess: &Cluster{
				Name:                     "testcluster",
				ApiEndpoint:              "https://api.testcluster.cluster.example.com:443",
				ControlPlaneEndpointHost: "api.testcluster.cluster.example.com",
				ControlPlaneEndpointPort: 443,
				Region:                   "us-east-1",
				ClusterGroup:             "test-clusters",
				Environment:              "test",
				CIDR:                     []string{"192.168.0.0/24"},
				ControlPlane:             &ClusterControlPlane{Provider: "kops"},
				Infrastructure:           &ClusterInfrastructure{Provider: "kops"},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "testcluster",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestCluster("testcluster2", "testcluster-kops-cp2", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster2", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*Cluster)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *Cluster", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := GetCluster(k, request.ResourceName)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetCluster_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetCluster should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Could not find cluster nonexistentcluster",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "nonexistentcluster",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "GetCluster should return Error for cluster without controlplane",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Cluster testcluster have an invalid configuration",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "testcluster",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "", "", "", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "GetCluster should return Error for cluster without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Cluster testcluster have an invalid configuration",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "testcluster",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "", "", ""),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			_, err := GetCluster(k, request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListClusters_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListClusters should return Success for one Cluster",
			ExpectedSuccess: []*Cluster{
				&Cluster{
					Name:                     "testcluster1",
					ApiEndpoint:              "https://api.testcluster1.cluster.example.com:443",
					ControlPlaneEndpointHost: "api.testcluster1.cluster.example.com",
					ControlPlaneEndpointPort: 443,
					Region:                   "us-east-1",
					ClusterGroup:             "test-clusters",
					Environment:              "test",
					CIDR:                     []string{"192.168.0.0/24"},
					ControlPlane:             &ClusterControlPlane{Provider: "kops"},
					Infrastructure:           &ClusterInfrastructure{Provider: "kops"},
				},
				&Cluster{
					Name:                     "testcluster2",
					ApiEndpoint:              "https://api.testcluster2.cluster.example.com:443",
					ControlPlaneEndpointHost: "api.testcluster2.cluster.example.com",
					ControlPlaneEndpointPort: 443,
					Region:                   "us-east-1",
					ClusterGroup:             "test-clusters",
					Environment:              "test",
					CIDR:                     []string{"192.168.0.0/24"},
					ControlPlane:             &ClusterControlPlane{Provider: "kops"},
					Infrastructure:           &ClusterInfrastructure{Provider: "kops"},
				},
			},
			ExpectedClientError: nil,
			Request:             &test.K8sRequest{},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster1", "testcluster-kops-cp1", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster1", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestCluster("testcluster2", "testcluster-kops-cp2", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster2", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		expectedInfra, ok := testCase.ExpectedSuccess.([]*Cluster)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to []Cluster", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := ListClusters(k)
			assert.NilError(t, err)
			assert.Equal(t, len(expectedInfra), len(response))
			for index, expectedCluster := range expectedInfra {
				response[index].ControlPlane = nil
				response[index].Infrastructure = nil
				expectedCluster.ControlPlane = nil
				expectedCluster.Infrastructure = nil
				assert.Assert(t, reflect.DeepEqual(expectedCluster, response[index]))

			}

		})
	}
}

func Test_ValidateClusterComponents_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "ValidateClusterComponents Should return success for a valid cluster",
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "testcluster",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		cluster, _ := k.GetCluster(request.ResourceName)
		t.Run(testCase.Name, func(t *testing.T) {
			err := ValidateClusterComponents(cluster)
			assert.NilError(t, err)
		})
	}
}

func Test_ValidateClusterComponents_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Should return an error for a cluster without controplane",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Cluster doesn't have a ControlPlane Reference",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "", "", "", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "Should return an error for a cluster without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Cluster doesn't have an infrastructure Reference",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "", "", ""),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)
		cluster, _ := testCase.K8sTestResources[0].(*clusterapiv1beta1.Cluster)
		t.Run(testCase.Name, func(t *testing.T) {
			err := ValidateClusterComponents(cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}
