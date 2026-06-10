package mail

import (
	"github.com/go-mail/mail/v2"
	"github.com/stringptr/SiGizi/backend/internal/config"
)

type Sender struct {
	dialer   *mail.Dialer
	from     string
	fromName string
}

func NewSender(cfg config.SMTPConfig) *Sender {
	d := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return &Sender{dialer: d, from: cfg.FromEmail, fromName: cfg.FromName}
}

func (s *Sender) Send(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetAddressHeader("From", s.from, s.fromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}
