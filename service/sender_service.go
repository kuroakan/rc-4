package service

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

type EmailService struct {
	client *mail.Client
	from   string
}

func NewEmailService(client *mail.Client, from string) *EmailService {
	return &EmailService{
		client: client,
		from:   from,
	}
}

func (s *EmailService) SendMail(email, text, subject string) error {
	message := mail.NewMsg()

	err := message.From(s.from)
	if err != nil {
		return fmt.Errorf("parse FROM address: %w", err)
	}

	if err := message.To(email); err != nil {
		return fmt.Errorf("parse TO address: %w", err)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, text)

	if err := s.client.DialAndSend(message); err != nil {
		return fmt.Errorf("send mail: %s", err)
	}

	return nil
}
