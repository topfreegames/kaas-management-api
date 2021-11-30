package test

import "github.com/topfreegames/kaas-management-api/util/clientError"

type TestCase struct {
	Name            string
	ExpectedSuccess interface{}
	ExpectedError   *clientError.ClientError
	Request         interface{}
}
