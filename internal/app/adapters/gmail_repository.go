package adapters

import (
	"fmt"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpServerAddr = "smtp.gmail.com"
	smtpPort       = 587
)

type GmailSender struct {
	name     string
	email    string
	password string
}

func NewGmailSender(name string, email string, password string) core.EmailRepository {
	return &GmailSender{
		name:     name,
		email:    email,
		password: password,
	}
}

func (g *GmailSender) SendEmail(to []string, subject string, body string) error {
	if env := config.EnvVars.ENV; env != "release" {
		return nil
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", g.name, g.email)
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	smtpAuth := smtp.PlainAuth("", g.email, g.password, smtpServerAddr)
	return e.Send(fmt.Sprintf("%s:%d", smtpServerAddr, smtpPort), smtpAuth)
}
