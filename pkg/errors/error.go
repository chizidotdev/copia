package errors

type Code int

const (
	ErrorBadRequest Code = iota
	ErrorNotFound
	ErrorInternal
	ErrorUnauthorized
	ErrorForbidden
)

type ErrResponse struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func Errorf(code Code, message string) *ErrResponse {
	return &ErrResponse{
		Code:    code,
		Message: message,
	}
}

func (e *ErrResponse) Error() string {
	return e.Message
}
