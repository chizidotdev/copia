package core


type EmailRepository interface {
	SendEmail(to []string, subject string, body string) error
}
