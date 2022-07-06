package k8s

import (
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapiexpv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"testing"
)

func Test_ValidateMachineTemplateComponents_Error(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:            "Should return an error for a NodeGroup without infrastructure",
			ExpectedSuccess: nil,
			ExpectedClientError: &clientError.ClientError{
				ErrorCause:           nil,
				ErrorDetailedMessage: "MachineTemplate doesn't have an infrastructure Reference",
				ErrorMessage:         clientError.InvalidConfiguration,
			},
			K8sTestResources: []runtime.Object{
				test.NewTestMachinePool("TestMachinePool", "TestCluster2", "", "", ""),
			},
		},
	}

	fakeClient := test.NewK8sFakeDynamicClient()
	k := &Kubernetes{K8sAuth: &Auth{
		DynamicClient: fakeClient,
	}}

	for _, testCase := range testCases {
		k.K8sAuth.DynamicClient = test.NewK8sFakeDynamicClientWithResources(testCase.K8sTestResources...)
		machinePool, _ := testCase.K8sTestResources[0].(*clusterapiexpv1beta1.MachinePool)
		t.Run(testCase.Name, func(t *testing.T) {
			err := ValidateMachineTemplateComponents(machinePool.Spec.Template)
			assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
			assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
		})
	}
}
