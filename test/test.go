package test

import (
	errorResponse "github.com/topfreegames/kaas-management-api/api/error"
	"github.com/topfreegames/kaas-management-api/util/clientError"
	"k8s.io/apimachinery/pkg/runtime"
)

type TestCase struct {
	Name                string
	ExpectedSuccess     interface{}
	ExpectedClientError *clientError.ClientError
	ExpectedHTTPError   *errorResponse.ClientErrorResponse
	TestK8sResources    []runtime.Object
	Request             interface{}
}
