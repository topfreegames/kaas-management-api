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
			Version:      "v9",
			EndpointName: "myendpoint",
			Path:         "/v9/myendpoint/",
		},
		"Request": map[string]string{
			"Version":          "v9",
			"EndpointNamePath": "myendpoint",
		},
	}

	request, _ := testCase["Request"].(map[string]string)
	expected, _ := testCase["ExpectedSuccess"].(*ApiEndpoint)

	t.Run(testCase["Name"].(string), func(t *testing.T) {
		endpoint := NewApiEndpoint(request["Version"], request["EndpointNamePath"])
		assert.Assert(t, reflect.DeepEqual(expected, endpoint))

	})
}
