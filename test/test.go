package test

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// Case Default template structure for table driven test
type Case struct {
	ExpectedBody interface{}
	ExpectedCode int
	Request      *Request
}

// Request Represents an Mock of HTTP request
type Request struct {
	Method string
	Body   io.Reader
	Path   string
}

// Run executes the Cases
func (r *Request) Run(handler http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(r.Method, r.Path, r.Body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}
