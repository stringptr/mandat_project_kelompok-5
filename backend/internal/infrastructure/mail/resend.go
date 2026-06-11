package mail

import (
	"github.com/resend/resend-go/v2"
	"github.com/stringptr/SiGizi/backend/internal/config"
)

type ResendSender struct {
	client   *resend.Client
	from     string
	fromName string
}

func NewResendSender(cfg config.MailConfig) *ResendSender {
	client := resend.NewClient(cfg.ResendAPIKey)
	return &ResendSender{client: client, from: cfg.FromEmail, fromName: cfg.FromName}
}

func (s *ResendSender) Send(to, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    s.fromName + " <" + s.from + ">",
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}
	_, err := s.client.Emails.Send(params)
	return err
}
