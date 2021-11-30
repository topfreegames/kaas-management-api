package test

import (
	"github.com/topfreegames/kaas-management-api/internal/k8s"
	"io"
	"k8s.io/client-go/dynamic/fake"
	"net/http"
	"net/http/httptest"
)

// HTTPTestCase Default template structure for table driven test
type HTTPTestCase struct {
	ExpectedBody interface{}
	ExpectedCode int
	Request      *HTTPTestRequest
}

// HTTPTestRequest Represents an Mock of HTTP request
type HTTPTestRequest struct {
	Method string
	Body   io.Reader
	Path   string
}

// RunHTTPTest executes the Cases
func (r *HTTPTestRequest) RunHTTPTest(handler http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(r.Method, r.Path, r.Body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func NewFakeKubernetesClient() *k8s.Kubernetes{
	fakeClient := fake.FakeDynamicClient{}
	k := &k8s.Kubernetes{K8sAuth: &k8s.Auth{
		DynamicClient: &fakeClient,
	}}
	return k
}