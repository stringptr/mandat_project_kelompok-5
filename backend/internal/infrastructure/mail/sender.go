package mail

import (
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/stringptr/SiGizi/backend/internal/config"
)

type Sender interface {
	Send(to, subject, body string) error
}

func New(cfg config.MailConfig) Sender {
	if cfg.Provider == "resend" {
		return NewResendSender(cfg)
	}
	return NewSMTPSender(cfg)
}

type SMTPSender struct {
	dialer   *mail.Dialer
	from     string
	fromName string
}

func NewSMTPSender(cfg config.MailConfig) *SMTPSender {
	d := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.SSL = cfg.SSL
	d.Timeout = 15 * time.Second
	return &SMTPSender{dialer: d, from: cfg.FromEmail, fromName: cfg.FromName}
}

func (s *SMTPSender) Send(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetAddressHeader("From", s.from, s.fromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}
