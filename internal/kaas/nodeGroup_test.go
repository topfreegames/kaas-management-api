package kaas

import (
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"reflect"
	"testing"
)

var TestKopsInfrastructure NodeInfrastructure = NodeInfrastructure{
	Name:        "TestMachinePool",
	Cluster:     "TestCluster1",
	Provider:    "Kops",
	Az:          []string{"us-east-1a"},
	MachineType: "m5.xlarge",
	Min:         nil,
	Max:         nil,
	Spec:        nil,
}

var TestDockerInfrastructure NodeInfrastructure = NodeInfrastructure{
	Name:     "docker",
	Provider: "docker",
	Cluster:  "docker-cluster",
	Az: []string{
		"local",
	},
	MachineType: "container",
	Spec:        nil,
}

func Test_GetNodeGroup_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetNodeGroup should return Success for MachinePool",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachinePool",
				Cluster:            "TestCluster1",
				InfrastructureName: "TestKopsMachinePool",
				InfrastructureKind: "KopsMachinePool",
				Infrastructure:     nil,
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
				test.NewTestKopsMachinePool("TestKopsMachinePool", "TestCluster1"),
			},
		},
		{
			Name: "GetNodeGroup should return Success for MachineDeployment",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachineDeployment",
				Cluster:            "TestCluster2",
				InfrastructureName: "TestDockerMachineTemplate",
				InfrastructureKind: "DockerMachineTemplate",
				Infrastructure:     nil,
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachineDeployment",
				Cluster:      "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
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
			response, err := GetNodeGroup(k, request.Cluster, request.ResourceName)
			response.Infrastructure = nil
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetNodeGroup_Error(t *testing.T) {
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
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "GetNodeGroup should return Error for node without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "NodeGroup TestMachinePool configuration is invalid",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "", "", ""),
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
			_, err := GetNodeGroup(k, request.Cluster, request.ResourceName)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_ListNodeGroup_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "ListNodeGroups should return Success for MachinePool",
			ExpectedSuccess: []*NodeGroup{
				{
					Name:               "TestMachinePool1",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool1",
					InfrastructureKind: "KopsMachinePool",
					Infrastructure:     nil,
					Replicas:           nil,
				},
				{
					Name:               "TestMachinePool2",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool2",
					InfrastructureKind: "KopsMachinePool",
					Infrastructure:     nil,
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
				test.NewTestKopsMachinePool("TestKopsMachinePool1", "TestCluster1"),
				test.NewTestKopsMachinePool("TestKopsMachinePool2", "TestCluster1"),
			},
		},
		{
			Name: "ListNodeGroups should return Success for MachineDeployment",
			ExpectedSuccess: []*NodeGroup{
				{
					Name:               "TestMachineDeployment1",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate1",
					InfrastructureKind: "DockerMachineTemplate",
					Infrastructure:     nil,
					Replicas:           nil,
				},
				{
					Name:               "TestMachineDeployment2",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate2",
					InfrastructureKind: "DockerMachineTemplate",
					Infrastructure:     nil,
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
				test.NewTestKopsMachinePool("TestKopsMachinePool", "TestCluster1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.([]*NodeGroup)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to []*NodeGroup", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := ListNodeGroups(k, request.Cluster)
			for _, ng := range response {
				ng.Infrastructure = nil
			}
			lenExpected := len(expectedInfra)
			lenResponse := len(response)
			assert.NilError(t, err)
			assert.Equal(t, lenExpected, lenResponse)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_ListNodeGroup_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "ListNodeGroups should return Error for non-existent resources in cluster",
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
			Name:            "ListNodeGroups should return Error for non-existent cluster",
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
		{
			Name:            "ListNodeGroups should return EmptyResponse for nodegroup without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "No valid NodeGroups were found in the cluster TestCluster2, some Nodegroups have invalid configuration",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "", "", ""),
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
			_, err := ListNodeGroups(k, request.Cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_GetNodeGroupConfig_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetNodeGroupConfig should return Success for MachinePool",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachinePool",
				Cluster:            "TestCluster1",
				InfrastructureName: "TestKopsMachinePool",
				InfrastructureKind: "KopsMachinePool",
				Infrastructure:     nil,
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name: "GetNodeGroupConfig should return Success for MachineDeployment",
			ExpectedSuccess: &NodeGroup{
				Name:               "TestMachineDeployment",
				Cluster:            "TestCluster2",
				InfrastructureName: "TestDockerMachineTemplate",
				InfrastructureKind: "DockerMachineTemplate",
				Infrastructure:     nil,
				Replicas:           nil,
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceName: "TestMachineDeployment",
				Cluster:      "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
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
			ng := &NodeGroup{
				Name:    request.ResourceName,
				Cluster: request.Cluster,
			}
			err := ng.getNodeGroupConfig(k)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, ng))
		})
	}
}

func Test_GetNodeGroupConfig_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetNodeGroupConfig should return Error for non-existent resource",
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
			Name:            "GetNodeGroupConfig should return Error for non-existent cluster",
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
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "DockerMachineTemplate", "TestDockerMachineTemplate", "infrastructure.cluster.x-k8s.io/v1beta1"),
			},
		},
		{
			Name:            "GetNodeGroupConfig should return Error for node without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "MachinePool TestMachinePool configuration is invalid",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			Request: &test.K8sRequest{
				ResourceName: "TestMachinePool",
				Cluster:      "TestCluster1",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestCluster1-TestMachinePool", "TestCluster1", "", "", ""),
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
			ng := &NodeGroup{
				Name:    request.ResourceName,
				Cluster: request.Cluster,
			}
			err := ng.getNodeGroupConfig(k)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}

func Test_GetNodeGroupListConfig_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetNodeGroupListConfig should return Success for MachinePool",
			ExpectedSuccess: []*NodeGroup{
				{
					Name:               "TestMachinePool1",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool1",
					InfrastructureKind: "KopsMachinePool",
					Infrastructure:     nil,
					Replicas:           nil,
				},
				{
					Name:               "TestMachinePool2",
					Cluster:            "TestCluster1",
					InfrastructureName: "TestKopsMachinePool2",
					InfrastructureKind: "KopsMachinePool",
					Infrastructure:     nil,
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
			Name: "GetNodeGroupListConfig should return Success for MachineDeployment",
			ExpectedSuccess: []*NodeGroup{
				{
					Name:               "TestMachineDeployment1",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate1",
					InfrastructureKind: "DockerMachineTemplate",
					Infrastructure:     nil,
					Replicas:           nil,
				},
				{
					Name:               "TestMachineDeployment2",
					Cluster:            "TestCluster2",
					InfrastructureName: "TestDockerMachineTemplate2",
					InfrastructureKind: "DockerMachineTemplate",
					Infrastructure:     nil,
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
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, ok := testCase.ExpectedSuccess.([]*NodeGroup)
		if !ok {
			log.Fatalf("Failed converting Success struct from test \"%s\" to []*NodeGroup", testCase.Name)
		}
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := GetNodeGroupListConfig(k, request.Cluster)
			lenExpected := len(expectedInfra)
			lenResponse := len(response)
			assert.NilError(t, err)
			assert.Equal(t, lenExpected, lenResponse)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetNodeGroupListConfig_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "GetNodeGroupListConfig should return Error for non-existent resources in cluster",
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
			Name:            "GetNodeGroupListConfig should return Error for non-existent cluster",
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
		{
			Name:            "ListNodeGroups should return EmptyResponse for nodegroup without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "No valid NodeGroups were found in the cluster TestCluster2, some Nodegroups have invalid configuration",
				ErrorMessage:         clientError.EmptyResponse,
			},
			Request: &test.K8sRequest{
				Cluster: "TestCluster2",
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster1", "KopsMachinePool", "TestKopsMachinePool", "infrastructure.cluster.x-k8s.io/v1alpha1"),
				test.NewTestMachineDeployment("TestCluster2-TestMachineDeployment", "TestCluster2", "", "", ""),
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
			_, err := GetNodeGroupListConfig(k, request.Cluster)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}
