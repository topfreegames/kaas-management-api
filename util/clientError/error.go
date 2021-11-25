package error

import (
	"github.com/gin-gonic/gin"
	apiError "github.com/topfreegames/kaas-management-api/api/error"
)

type ClientError struct {
	ErrorCause   error
	ErrorMessage string
	ErrorCode    int // TODO define error codes
}

func (e ClientError) Error() string {
	if e.ErrorCause == nil {
		return e.ErrorMessage
	}
	return e.ErrorMessage + " : " + e.ErrorCause.Error()
}

func NewClientError(errorCause error, errorMessage string) error {
	clientError := &ClientError{
		ErrorCause:   errorCause,
		ErrorMessage: errorMessage,
	}
	return clientError
}

func ErrorHandler(c *gin.Context, err error, errorMessage string, httpCode int) {
	clientErrorResponse := &apiError.ClientErrorResponse{
		ErrorMessage: errorMessage,
	}
	c.JSON(httpCode, clientErrorResponse)
}
