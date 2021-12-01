package test

import (
	"github.com/gin-gonic/gin"
	"github.com/topfreegames/kaas-management-api/api"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// HTTPTestExpectedResponse
type HTTPTestExpectedResponse struct {
	ExpectedBody interface{}
	ExpectedCode int
}

// HTTPTestRequest Represents an Mock of HTTP request
type HTTPTestRequest struct {
	Method string
	Body   io.Reader
	Path   string
}

// GetK8sRequest returns the request of the test as an instance of the struct *HTTPTestRequest
func (t TestCase) GetHTTPRequest() *HTTPTestRequest {
	request, ok := t.Request.(*HTTPTestRequest)
	if !ok {
		log.Fatalf("Could not convert TestCase %s Request to HTTPTestRequest", t.Name)
	}
	return request
}

func SetupEndpointRouter(endpoint *api.ApiEndpoint) *api.ApiEndpoint {
	endpoint.Router = nil
	endpoint.RouterGroup = nil
	endpoint.CreatePublicRouterGroup(gin.Default())
	return endpoint
}

// RunHTTPTest executes the Cases
func (r *HTTPTestRequest) RunHTTPTest(handler http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(r.Method, r.Path, r.Body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}
