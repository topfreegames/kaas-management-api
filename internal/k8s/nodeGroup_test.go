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
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"testing"
)

func Test_GetNodeGroup_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetNodeGroup should return Success for MachinePool",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachinePool",
				Cluster:            "TestCluster1",
				InfrastructureName: "TestKopsMachinePool",
				InfrastructureKind: "KopsMachinePool",
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name: "GetNodeGroup should return Success for MachineDeployment",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachineDeployment",
				Cluster:            "TestCluster2",
				InfrastructureName: "TestDockerMachineTemplate",
				InfrastructureKind: "DockerMachineTemplate",
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachineDeployment",
				Cluster:      "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
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
		expectedInfra, ok := testCase.ExpectedSuccess.(*NodeGroup)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *NodeGroup", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.GetNodeGroup(request.Cluster, request.ResourceName)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetNodeGroup_ErrorNotFound(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetNodeGroup should return Error for non-existent resource",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Could not find the NodeGroup nonexistent in the cluster TestCluster1",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "nonexistent",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "GetNodeGroup should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "Could not find the NodeGroup TestMachinePool in the cluster TestCluster3",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
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
			_, err := k.GetNodeGroup(request.Cluster, request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListNodeGroup_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListNodeGroup should return Success for MachinePool",
			ExpectedSuccess: []NodeGroup{
				{
					Name:               "TestMachinePool1",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool1",
					InfrastructureKind: "KopsMachinePool",
					Replicas:           nil,
				},
				{
					Name:               "TestMachinePool2",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool2",
					InfrastructureKind: "KopsMachinePool",
					Replicas:           nil,
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool1", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("TestMachinePool2", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool2", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name: "ListNodeGroup should return Success for MachineDeployment",
			ExpectedSuccess: []NodeGroup{
				{
					Name:               "TestMachineDeployment1",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate1",
					InfrastructureKind: "DockerMachineTemplate",
					Replicas:           nil,
				},
				{
					Name:               "TestMachineDeployment2",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate2",
					InfrastructureKind: "DockerMachineTemplate",
					Replicas:           nil,
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment1", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate1", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachineDeployment("TestMachineDeployment2", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate2", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.([]NodeGroup)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to []NodeGroup", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.ListNodeGroup(request.Cluster)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_ListNodeGroup_ErrorEmptyResponse(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListNodeGroup should return Error for non-existent resources in cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "No NodeGroups were found in the cluster TestCluster1",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestKopsMachinePool("test", "TestCluster1"),
				test.NewTestMachinePool("TestMachinePool", "TestCluster2", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "ListNodeGroup should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "No NodeGroups were found in the cluster TestCluster3",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
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
			_, err := k.ListNodeGroup(request.Cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_GetMachinePool_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "GetMachinePool should return Success for MachinePool",
			ExpectedSuccess:     test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachinePool("TestMachinePool2", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
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

func Test_GetMachinePool_ErrorNotFound(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetMachinepool should return Error for non-existent resource",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested machinepool nonexistent was not found for the cluster TestCluster1!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "nonexistent",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "GetMachinepool should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested machinepool TestMachinePool was not found for the cluster TestCluster3!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
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

func Test_ListMachinePool_ErrorEmptyResponse(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListMachinePool should return Error for non-existent resources in cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no Machinepools were found for the cluster TestCluster1!",
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
				ErrorDetailedMessage: "no Machinepools were found for the cluster TestCluster3!",
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

func Test_GetMachineDeployment_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "GetMachineDeployment should return Success for MachineDeployment",
			ExpectedSuccess:     test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachineDeployment",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachineDeployment("TestMachineDeployment2", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiv1beta1.MachineDeployment)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *clusterapiv1beta1.MachineDeployment", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.GetMachineDeployment(request.Cluster, request.ResourceName)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetMachineDeployment_ErrorNotFound(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetMachineDeployment should return Error for non-existent resource",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested machinedeployment nonexistent was not found for the cluster TestCluster1!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "nonexistent",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "GetMachinepool should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "The requested machinedeployment TestMachineDeployment was not found for the cluster TestCluster3!",
				ErrorMessage:         clientError.ResourceNotFound,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachineDeployment",
				Cluster:      "TestCluster3",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
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
			_, err := k.GetMachineDeployment(request.Cluster, request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListMachineDeployment_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListMachineDeployment should return Success for two MachineDeployments",
			ExpectedSuccess: &clusterapiv1beta1.MachineDeploymentList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineDeploymentList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiv1beta1.MachineDeployment{
					*test.NewTestMachineDeployment("TestMachineDeployment1", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
					*test.NewTestMachineDeployment("TestMachineDeployment2", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment1", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachineDeployment("TestMachineDeployment2", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name: "ListMachineDeployment should return Success for one MachineDeployment",
			ExpectedSuccess: &clusterapiv1beta1.MachineDeploymentList{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineDeploymentList",
					APIVersion: "cluster.x-k8s.io/v1beta1",
				},
				ListMeta: metav1.ListMeta{},
				Items: []clusterapiv1beta1.MachineDeployment{
					*test.NewTestMachineDeployment("TestMachineDeployment3", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				},
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				Cluster: "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment3", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.(*clusterapiv1beta1.MachineDeploymentList)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to *clusterapiv1beta1.MachineDeploymentList", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := k.ListMachineDeployment(request.Cluster)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_ListMachineDeployment_ErrorEmptyResponse(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListMachineDeployment should return Error for non-existent resources in cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no Machinedeployments were found for the cluster TestCluster1!",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool1", "infrastructure.cluster.x-k8s.io/v1alpha1"),
			},
		},
		{
			Name:            "ListMachineDeployment should return Error for non-existent cluster",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "no Machinedeployments were found for the cluster TestCluster3!",
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
			_, err := k.ListMachineDeployment(request.Cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}
