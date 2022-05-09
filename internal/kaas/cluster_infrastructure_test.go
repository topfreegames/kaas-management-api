package kaas

import (
	"github.com/topfreegames/kaas-management-api/test"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_GetClusterInfrastructure_Success(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "GetClusterInfrastructure should return Success for Kops",
			ExpectedSuccess: &ClusterInfrastructure{
				Provider: "kops",
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceKind: "KopsAWSCluster",
			},
		},
		{
			Name: "GetClusterInfrastructure should return Success for Docker",
			ExpectedSuccess: &ClusterInfrastructure{
				Provider: "docker",
			},
			ExpectedClientError: nil,
			Request: &test.K8sRequest{
				ResourceKind: "DockerCluster",
			},
		},
	}

	for _, testCase := range testCases {
		request := testCase.GetK8sRequest()
		expectedInfra, _ := testCase.ExpectedSuccess.(*ClusterInfrastructure)

		t.Run(testCase.Name, func(t *testing.T) {
			response, err := GetClusterInfrastructure(request.ResourceKind)
			assert.NilError(t, err)
			assert.Assert(t, reflect.DeepEqual(expectedInfra, response))
		})
	}
}

func Test_GetClusterInfrastructure_ErrorKindNotFound(t *testing.T) {
	testCase := test.TestCase{
		Name:            "GetClusterInfrastructure should return Kind not found",
		ExpectedSuccess: nil,
		ExpectedClientError: &clientError.ClientError{
			ErrorCause:           nil,
			ErrorDetailedMessage: "The Kind NonExistentKind could not be found",
			ErrorMessage:         clientError.KindNotFound,
		},
		Request: &test.K8sRequest{
			ResourceKind: "NonExistentKind",
		},
	}

	request := testCase.GetK8sRequest()

	t.Run(testCase.Name, func(t *testing.T) {
		_, err := GetControlPlane(request.ResourceKind)
		assert.ErrorContains(t, err, testCase.ExpectedClientError.Error())
		assert.Assert(t, test.AssertClientError(err, testCase.ExpectedClientError))
	})
}
