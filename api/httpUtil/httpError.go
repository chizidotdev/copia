package httpUtil

type Code int

const (
	ErrorBadRequest Code = iota
	ErrorNotFound
	ErrorInternal
	ErrorUnauthorized
	ErrorForbidden
)

type HttpError struct {
	Code      Code   `json:"code"`
	MessageID string `json:"messageID"`
	Message   string `json:"message"`
	Reason    string `json:"reason"`
}

func Errorf(errResp HttpError) *HttpError {
	return &HttpError{
		Code:      errResp.Code,
		MessageID: errResp.MessageID,
		Message:   errResp.Message,
		Reason:    errResp.Reason,
	}
}

func (e *HttpError) Error() string {
	return e.Message
}
