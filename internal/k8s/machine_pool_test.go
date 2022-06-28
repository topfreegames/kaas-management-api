package k8s

import (
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"reflect"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"testing"
)

func Test_GetMachinePool_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "GetMachinePool should return Success for MachinePool",
			ExpectedSuccess:     test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestCluster1-TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("TestCluster1-TestMachinePool2", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiexpv1beta1.MachinePool)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *v1beta1.MachinePool", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.GetMachinePool(request.Cluster, request.ResourceName)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetMachinePool_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetMachinePool should return Error for non-existent resource",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested MachinePool nonexistent was not found for the cluster TestCluster1!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "nonexistent",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "GetMachinePool should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested MachinePool TestMachinePool was not found for the cluster TestCluster3!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "GetMachinePool should return Error for a MachinePool without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "MachinePool TestCluster1-TestMachinePool doesn't have a valid configuration",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestCluster1-TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "", "", ""),
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
			_, err := k.GetMachinePool(request.Cluster, request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListMachinePool_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListMachinePool should return Success for two MachinePools",
			ExpectedSuccess: &clusterapiexpv1beta1.MachinePoolList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachinePoolList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiexpv1beta1.MachinePool{
					*test.NewTestMachinePool("TestMachinePool1", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
					*test.NewTestMachinePool("TestMachinePool2", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool2", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool1", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("TestMachinePool2", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool2", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name: "ListMachinePool should return Success for one MachinePool",
			ExpectedSuccess: &clusterapiexpv1beta1.MachinePoolList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachinePoolList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiexpv1beta1.MachinePool{
					*test.NewTestMachinePool("TestMachinePool3", "TestCluster2", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool3", "TestCluster2", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiexpv1beta1.MachinePoolList)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *clusterapiexpv1beta1.MachinePoolList", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.ListMachinePool(request.Cluster)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_ListMachinePool_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListMachinePool should return Error for non-existent resources in cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no MachinePools were found for the cluster TestCluster1!",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachinePool("TestMachinePool", "TestCluster2", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "ListMachinePool should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no MachinePools were found for the cluster TestCluster3!",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster2", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
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
			_, err := k.ListMachinePool(request.Cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}
