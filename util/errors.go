package util

type ErrorCode int

const (
	ErrorBadRequest ErrorCode = iota
	ErrorNotFound
	ErrorInternal
	ErrorUnauthorized
	ErrorForbidden
)

type ErrResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func Errorf(code ErrorCode, message string) *ErrResponse {
	return &ErrResponse{
		Code:    code,
		Message: message,
	}
}

func (e *ErrResponse) Error() string {
	return e.Message
}

// ErrorMessage
//
// Deprecated: use local package errorMessage instead.
func ErrorMessage(message string) string {
	return message
}
