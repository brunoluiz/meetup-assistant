package email

import (
	"context"
	"time"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/mailgun/mailgun-go/v4"
)

type Mailgun struct {
	domain  string
	mailgun mailgun.Mailgun
}

func NewMailgun(domain, token string) *Mailgun {
	return &Mailgun{
		domain:  domain,
		mailgun: mailgun.NewMailgun(domain, token),
	}
}

func (m *Mailgun) Send(ctx context.Context, target channel.Target, subject, body string) error {
	sender := "contact@" + m.domain

	message := m.mailgun.NewMessage(sender, subject, "", target.Address)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err := m.mailgun.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
