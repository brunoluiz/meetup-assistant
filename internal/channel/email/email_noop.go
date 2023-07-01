package email

import (
	"context"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
)

type Noop struct {
}

func NewNoop() *Noop {
	return &Noop{}
}

func (m *Noop) Send(ctx context.Context, target channel.Target, subject, body string) error {
	return nil
}
