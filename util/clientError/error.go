package clientError

import (
	"github.com/gin-gonic/gin"
	apiError "github.com/topfreegames/kaas-management-api/api/error"
)

type ClientError struct {
	ErrorCause           error
	ErrorDetailedMessage string
	ErrorMessage         string
	ErrorCode            int // TODO define clientError codes
}

func (e ClientError) Error() string {
	if e.ErrorCause == nil {
		return e.ErrorMessage + ": " + e.ErrorDetailedMessage
	}
	return e.ErrorMessage + ": " + e.ErrorDetailedMessage + " caused by: " + e.ErrorCause.Error()
}

func NewClientError(errorCause error, errorMessage string, errorDetailedMessage string) error {
	clientError := &ClientError{
		ErrorCause:           errorCause,
		ErrorMessage:         errorMessage,
		ErrorDetailedMessage: errorDetailedMessage,
	}
	return clientError
}

func ErrorHandler(c *gin.Context, err error, errorMessage string, httpCode int) {
	var errorMsg string

	clientErr, ok := err.(*ClientError)
	if ok {
		errorMsg = clientErr.ErrorMessage
	} else {
		errorMsg = UnexpectedError
	}

	clientErrorResponse := &apiError.ClientErrorResponse{
		ErrorMessage: errorMessage,
		ErrorType:    errorMsg,
		HttpCode:     httpCode,
	}
	c.JSON(httpCode, clientErrorResponse)
}
