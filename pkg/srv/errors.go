package srv

import (
	"fmt"
	"net/http"
)

// ErrMsg describes a server error
type ErrMsg struct {
	Kind   string      `json:"kind"`
	Status int         `json:"status"`
	Msg    string      `json:"msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// NewErr creates an ErrMsg with given msg
func NewErr(kind string, status int, msg string) ErrMsg {
	return ErrMsg{kind, status, msg, nil}
}

// NewErrWithData creates an ErrMsg including extra data. Make sure data is
// json encodable
func NewErrWithData(kind string, status int, msg string, data interface{}) ErrMsg {
	return ErrMsg{kind, status, msg, data}
}

// NewErrInternal creates an ErrMsg for internal server errors
func NewErrInternal() ErrMsg {
	return ErrMsg{"ERR_INTERNAL", http.StatusInternalServerError, "Unexpected error happened", nil}
}

// NewErrMissingParam creates an ErrMsg for a missing parameter
func NewErrMissingParam(param string) ErrMsg {
	msg := fmt.Sprintf("Query parameter '%s' is missing", param)
	return NewErr("ERR_MISSING_PARAM", http.StatusBadRequest, msg)
}

// NewErrBadValue creates an ErrMsg for a bad value
func NewErrBadValue(param, value string) ErrMsg {
	msg := fmt.Sprintf("Bad value '%s' for param '%s'", param, value)
	return NewErr("ERR_BAD_VALUE", http.StatusBadRequest, msg)
}

// NewErrBadRequest creates an ErrMsg for bad requests describing what are the
// mistakes
func NewErrBadRequest(mistakes map[string]interface{}) ErrMsg {
	return NewErrWithData("ERR_BAD_REQUEST", http.StatusBadRequest, "Check data for invalid parameters", mistakes)
}

// NewErrUnsupported creates an ErrMsg for an unsupported feature
func NewErrUnsupported(feature string, vals ...interface{}) ErrMsg {
	var msg string
	if len(vals) > 0 {
		msg = fmt.Sprintf("%s is not supported", fmt.Sprintf(feature, vals))
	} else {
		msg = fmt.Sprintf("%s is not supported", feature)
	}

	return NewErr("ERR_UNSUPPORTED", http.StatusBadRequest, msg)
}

// Error reports error message
func (e ErrMsg) Error() string {
	return e.Msg
}
