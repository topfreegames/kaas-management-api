package test

import (
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

func Param(param string) string {
	return ":" + param + "/"
}

func Path(path string) string {
	return path + "/"
}

// RunHTTPTest executes the Cases
func (r *HTTPTestRequest) RunHTTPTest(handler http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(r.Method, r.Path, r.Body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}
