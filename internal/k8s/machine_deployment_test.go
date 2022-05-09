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

func Test_GetMachineDeployment_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:                "GetMachineDeployment should return Success for MachineDeployment",
			ExpectedSuccess:     test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestCluster1-TestMachineDeployment",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment2", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
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

func Test_GetMachineDeployment_Error(t *testing.T) {
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
				test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
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
				test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment", "TestCluster1", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "GetMachinepool should return Error for a machinePool without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "MachineDeployment TestCluster1-TestMachineDeployment doesn't have a valid configuration",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestCluster1-TestMachineDeployment",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachineDeployment("TestCluster1-TestMachineDeployment", "TestCluster1", "", "", ""),
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

func Test_ListMachineDeployment_Error(t *testing.T) {
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
