package email

import (
	"context"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/davecgh/go-spew/spew"
)

type Noop struct {
	Template    Template
	Idempotency Idempotency
}

func NewNoop(template Template, idempotency Idempotency) *Noop {
	return &Noop{
		Template:    template,
		Idempotency: idempotency,
	}
}

func (m *Noop) Send(
	ctx context.Context,
	templatePath string,
	target channel.Target,
	params map[string]any,
) error {
	content, err := m.Template.Render(ctx, templatePath, params)
	if err != nil {
		return err
	}

	spew.Dump(content)

	return nil
}

func (m *Noop) Type() channel.Type {
	return channel.TypeEmail
}
