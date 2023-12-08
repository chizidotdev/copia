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
	Code      Code   `json:"code"`
	MessageID string `json:"messageID"`
	Message   string `json:"message"`
	Reason    string `json:"reason"`
}

func Errorf(errResp ErrResponse) *ErrResponse {
	return &ErrResponse{
		Code:      errResp.Code,
		MessageID: errResp.MessageID,
		Message:   errResp.Message,
		Reason:    errResp.Reason,
	}
}

func (e *ErrResponse) Error() string {
	return e.Message
}
