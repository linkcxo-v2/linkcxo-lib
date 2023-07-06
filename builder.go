package linkcxo

type APIResponse struct {
	Data   any    `json:"data"`
	Status int    `json:"status"`
	Error  *Error `json:"error,omitempty"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type errorCode struct {
	Common commonErrorCode
}
type commonErrorCode struct {
	StatusUnsupportedMediaType     string
	StatusForbidden                string
	StatusNotFound                 string
	StatusUnauthorized             string
	StatusInternalServerError      string
	StatusUnprocessableEntityError string
}

var ErrorCode = errorCode{
	Common: commonErrorCode{
		StatusUnsupportedMediaType:     "AAS-00100",
		StatusForbidden:                "AAS-00101",
		StatusNotFound:                 "AAS-00102",
		StatusUnauthorized:             "AAS-00103",
		StatusUnprocessableEntityError: "AAS-00104",
	},
}

type ResponseBuilder struct {
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}
func (r *ResponseBuilder) BuildError(err error, msg string, httpError int) APIResponse {
	return APIResponse{
		Status: 0,
		Error: &Error{
			Message: msg,
			Code:    httpError,
		},
	}
}
func (r *ResponseBuilder) BuildSuccess(data any) APIResponse {

	return APIResponse{
		Data:   data,
		Status: 1,
	}
}
