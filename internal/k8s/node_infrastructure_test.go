package k8s

import (
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	clusterapikopsv1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kops/pkg/apis/kops/v1alpha2"
	"reflect"
	"testing"
)

func Test_GetNodeInfrastructure_Success(t *testing.T) {
	testCase := test.TestCase{
		Name: "GetNodeInfrastructure should return Success for a KopsMachinePool",
		ExpectedSuccess: &NodeInfrastructure{
			Name:        "kops-test",
			Cluster:     "test-cluster",
			Provider:    "kops",
			Az:          []string{"us-east-1a"},
			MachineType: "m5.xlarge",
			Min:         nil,
			Max:         nil,
			Spec: clusterapikopsv1alpha1.KopsMachinePoolSpec{
				KopsInstanceGroupSpec: v1alpha2.InstanceGroupSpec{
					MinSize:     nil,
					MaxSize:     nil,
					MachineType: "m5.xlarge",
					Subnets:     []string{"us-east-1a"},
				}},
		},
		ExpectedError: nil,
		Request: &test.K8sRequest{
			ResourceName: "kops-test",
			ResourceKind: "KopsMachinePool",
			Cluster:      "test-cluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("kops-test", "test-cluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()
	expectedInfra, _ := testCase.ExpectedSuccess.(*NodeInfrastructure)

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		response, err := k.GetNodeInfrastructure(request.Cluster, request.ResourceKind, request.ResourceName)
		assert.NilError(t, err)
		assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
	})
}

func Test_GetNodeInfrastructure_ErrorResourceNotFound(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetNodeInfrastructure should return resource not found error",
		ExpectedSuccess: nil,
		ExpectedError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "Could not retrieve the infrastructure",
			ErrorMessage:         "RESOURCE_NOT_FOUND",
		},
		Request: &test.K8sRequest{
			ResourceName: "NonExistentResource",
			ResourceKind: "KopsMachinePool",
			Cluster:      "test-cluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("OtherResource", "test-cluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := k.GetNodeInfrastructure(request.Cluster, request.ResourceKind, request.ResourceName)
		assert.ErrorContains(t, err, testCase.ExpectedError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedError))
	})
}

func Test_GetNodeInfrastructure_ErrorResourceNotFoundForCluster(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetNodeInfrastructure should return resource not found error for non-existent cluster",
		ExpectedSuccess: nil,
		ExpectedError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "Could not retrieve the infrastructure",
			ErrorMessage:         "RESOURCE_NOT_FOUND",
		},
		Request: &test.K8sRequest{
			ResourceName: "test-kops",
			ResourceKind: "KopsMachinePool",
			Cluster:      "nonExistentCluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("test-kops", "otherCluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := k.GetNodeInfrastructure(request.Cluster, request.ResourceKind, request.ResourceName)
		assert.ErrorContains(t, err, testCase.ExpectedError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedError))
	})
}

func Test_GetNodeInfrastructure_ErrorKindNotFound(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetNodeInfrastructure should return Kind not found",
		ExpectedSuccess: nil,
		ExpectedError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "The Kind NonExistentKind could not be found",
			ErrorMessage:         "KIND_NOT_FOUND",
		},
		Request: &test.K8sRequest{
			ResourceName: "kops-test",
			ResourceKind: "NonExistentKind",
			Cluster:      "test-cluster",
		},
	}

	request := testCase.GetK8sRequest()

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := k.GetNodeInfrastructure(request.Cluster, request.ResourceKind, request.ResourceName)
		assert.ErrorContains(t, err, testCase.ExpectedError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedError))
	})
}

func Test_GetKopsMachinePool_Success(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetKopsMachinePool should return Success for a KopsMachinePool",
		ExpectedSuccess: test.NewTestKopsMachinePool("myMachinePool", "mycluster"),
		ExpectedError:   nil,
		Request: &test.K8sRequest{
			ResourceName: "myMachinePool",
			Cluster:      "mycluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("myMachinePool", "mycluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()
	expectedInfra, _ := testCase.ExpectedSuccess.(*clusterapikopsv1alpha1.KopsMachinePool)

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		response, err := k.GetKopsMachinePool(request.Cluster, request.ResourceName)
		assert.NilError(t, err)
		assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
	})
}

func Test_GetKopsMachinePool_ErrorResourceNotFound(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetKopsMachinePool should return Error for a non-existent machinePool",
		ExpectedSuccess: nil,
		ExpectedError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "The requested KopsMachinePool NonExistentMachinePool was not found in namespace mycluster!",
			ErrorMessage:         "RESOURCE_NOT_FOUND",
		},
		Request: &test.K8sRequest{
			ResourceName: "NonExistentMachinePool",
			Cluster:      "mycluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("OtherResource", "mycluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := k.GetKopsMachinePool(request.Cluster, request.ResourceName)
		assert.ErrorContains(t, err, testCase.ExpectedError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedError))
	})
}

func Test_GetKopsMachinePool_ErrorResourceNotFoundNonExistentCluster(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetKopsMachinePool should return Error for a non-existent cluster",
		ExpectedSuccess: nil,
		ExpectedError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "The requested KopsMachinePool ExistentResource was not found in namespace NonExistentCluster!",
			ErrorMessage:         "RESOURCE_NOT_FOUND",
		},
		Request: &test.K8sRequest{
			ResourceName: "ExistentResource",
			Cluster:      "NonExistentCluster",
			TestResources: []runtime.Object{
				test.NewTestKopsMachinePool("ExistentResource", "ExistentCluster"),
			},
		},
	}

	request := testCase.GetK8sRequest()

	fakeClient := test.NewK8sFakeDynamicClientWithResources(request.TestResources...)
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := k.GetKopsMachinePool(request.Cluster, request.ResourceName)
		assert.ErrorContains(t, err, testCase.ExpectedError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedError))
	})
}
