package healthCheck

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	healthCheckv1 "github.com/topfreegames/kaas-management-api/apis/healthCheck"
	"github.com/topfreegames/kaas-management-api/test"
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	tests := map[string]test.Case{
		"Success health check response": {
			ExpectedCode: http.StatusOK,
			ExpectedBody: healthCheckv1.HealthCheck{
				Healthy: true,
			},
			Request: &test.Request{
				Method: "GET",
				Body:   nil,
				Path:   "/healthcheck",
			},
		},
	}

	router := gin.Default()
	router.GET("/healthcheck", HealthCheckHandler)
	for testMsg, testCase := range tests {
		t.Run(testMsg, func(t *testing.T) {
			w := testCase.Request.Run(router)
			assert.Equal(t, testCase.ExpectedCode, w.Code)
			expected, err := json.Marshal(testCase.ExpectedBody)
			assert.Nil(t, err)
			assert.Equal(t, string(expected), w.Body.String())
		})
	}
}