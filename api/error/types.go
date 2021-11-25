package error

type ClientErrorResponse struct {
	ErrorMessage string `json:"errormessage"`
	ErrorCode    int    `json:"errorcode,omitempty"`
	ErrorType    string `json:"errortype,omitempty"`
}
