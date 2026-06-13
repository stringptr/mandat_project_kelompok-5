package testutils

import "sync"

type CapturedEmail struct {
	To      string
	Subject string
	Body    string
}

type CaptureSender struct {
	mu     sync.Mutex
	emails []CapturedEmail
}

func NewCaptureSender() *CaptureSender {
	return &CaptureSender{}
}

func (s *CaptureSender) Send(to, subject, body string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails = append(s.emails, CapturedEmail{To: to, Subject: subject, Body: body})
	return nil
}

func (s *CaptureSender) SentEmails() []CapturedEmail {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]CapturedEmail, len(s.emails))
	copy(result, s.emails)
	return result
}

func (s *CaptureSender) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails = nil
}
