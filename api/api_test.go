package api

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_NewApiEndpoint_Success(t *testing.T) {
	testCase := map[string]interface{}{
		"Name": "NewApiEndpoint should create a new ApiEndpoint object",
		"ExpectedSuccess": &ApiEndpoint{
			version:      "v9",
			endpoint:     "myendpoint",
			EndpointPath: "/v9/myendpoint",
		},
		"Request": map[string]string{
			"version":  "v9",
			"endpoint": "myendpoint",
		},
	}

	request, _ := testCase["Request"].(map[string]string)
	expected, _ := testCase["ExpectedSuccess"].(*ApiEndpoint)

	t.Run(testCase["Name"].(string), func(t *testing.T) {
		endpoint := NewApiEndpoint(request["version"], request["endpoint"])
		assert.Assert(t, reflect.DeepEqual(expected, endpoint))

	})
}
