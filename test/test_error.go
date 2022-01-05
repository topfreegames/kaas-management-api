package test

import "github.com/topfreegames/kaas-management-api/util/clientError"

// AssertClientError Returns true if the Error is from the type *ClientError and if it is correctly set
func AssertClientError(err error, expectedErr *clientError.ClientError) bool {
	clientErr, ok := err.(*clientError.ClientError)
	if !ok {
		return false
	}
	if expectedErr.ErrorMessage != clientErr.ErrorMessage {
		return false
	}
	if expectedErr.ErrorDetailedMessage != clientErr.ErrorDetailedMessage {
		return false
	}
	return true
}
