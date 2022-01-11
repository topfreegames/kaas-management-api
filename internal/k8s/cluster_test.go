package k8s

import (
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"reflect"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"testing"
)

func Test_GetCluster_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "GetCluster should return Success for one cluster",
			ExpectedSuccess:     test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
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
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiv1beta1.Cluster)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *v1beta1.MachinePool", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.GetCluster(request.ResourceName)
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
				ErrorDetailedMessage: "The requested cluster nonexistentcluster was not found in namespace kubernetes-nonexistentcluster!",
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
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			_, err := k.GetCluster(request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListClusters_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListClusters should return Success for one Cluster",
			ExpectedSuccess: &clusterapiv1beta1.ClusterList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiv1beta1.Cluster{
					*test.NewTestCluster("testcluster1", "testcluster-kops-cp1", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster1", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				},
			},
			ExpectedClientError: nil,
			Request:             &test.K8sRequest{},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster1", "testcluster-kops-cp1", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster1", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name: "ListMachinePool should return Success for one MachinePool",
			ExpectedSuccess: &clusterapiv1beta1.ClusterList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiv1beta1.Cluster{
					*test.NewTestCluster("testcluster2", "testcluster-kops-cp2", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster2", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
					*test.NewTestCluster("testcluster3", "testcluster-kops-cp3", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster3", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				},
			},
			ExpectedClientError: nil,
			Request:             &test.K8sRequest{},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster2", "testcluster-kops-cp2", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster2", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
				test.NewTestCluster("testcluster3", "testcluster-kops-cp3", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "kops-cluster3", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiv1beta1.ClusterList)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *clusterapiexpv1beta1.MachinePoolList", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.ListClusters()
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_ListClusters_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListClusters should return EmptyResponse when have only one cluster without infrastructure reference",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no valid clusters were found, some clusters have invalid configuration",
				ErrorMessage:         clientError.EmptyResponse,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "", "", ""),
			},
		},
		{
			Name:            "ListClusters should return EmptyResponse when have only one cluster without control plane reference",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no valid clusters were found, some clusters have invalid configuration",
				ErrorMessage:         clientError.EmptyResponse,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "", "", "", "kops-cluster", "KopsAWSCluster", "controlplane.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			_, err := k.ListClusters()
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
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
	k := &Kubernetes{K8sAuth: &Auth{
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
				ErrorDetailedMessage: "Cluster doesn't have a infrastructure Reference",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestCluster("testcluster", "testcluster-kops-cp", "KopsControlPlane", "controlplane.cluster.x-k8s.io/v1alpha1", "", "", ""),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
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
