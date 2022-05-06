package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	healthCheckv1 "github.com/topfreegames/kaas-management-api/api/healthCheck"
	"github.com/topfreegames/kaas-management-api/test"
)

func TestHealthCheckHandler(t *testing.T) {
	testCase := test.TestCase{
		Name: "Healthcheck should return ok",
		ExpectedSuccess: test.HTTPTestExpectedResponse{
			ExpectedBody: healthCheckv1.HealthCheck{
				Healthy: true,
			},
			ExpectedCode: http.StatusOK,
		},
		ExpectedHTTPError: nil,
		Request: &test.HTTPTestRequest{
			Method: http.MethodGet,
			Body:   nil,
			Path:   healthCheckv1.Endpoint.Path,
		},
	}

	router := gin.Default()
	router.Handle(http.MethodGet, healthCheckv1.Endpoint.Path, HealthCheckHandler)

	request := testCase.GetHTTPRequest()
	expectedResponse, ok := testCase.ExpectedSuccess.(test.HTTPTestExpectedResponse)
	if !ok {
		log.Fatalf("Failed converting Success struct from test \"%s\" to *test.HTTPTestExpectedResponse", testCase.Name)
	}

	t.Run(testCase.Name, func(t *testing.T) {

		w := request.RunHTTPTest(router)

		assert.Equal(t, expectedResponse.ExpectedCode, w.Code)
		expected, err := json.Marshal(expectedResponse.ExpectedBody)
		assert.Nil(t, err)
		assert.Equal(t, string(expected), w.Body.String())
	})
}
