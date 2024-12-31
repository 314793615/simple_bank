package mail

import (
	"crypto/internal/nistec"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAddress = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)


type EmailSender interface {
	SendEmail() error
}

type GmailSender struct {
	name string
	fromEmailAddress string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) *GmailSender {
	return &GmailSender{
		name: name,
		fromEmailAddress: fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}


func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	bcc []string,
	cc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>",sender.name, sender.fromEmailAddress)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.HTML = []byte(content)
	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}	
	smtpAtuth := smtp.PlainAuth("", sender.name, sender.fromEmailPassword, smtpAddress)
	return e.Send(smtpServerAddress, smtpAtuth)
}