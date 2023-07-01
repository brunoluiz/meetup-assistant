package email

import (
	"context"
	"fmt"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/brunoluiz/meetup-assistant/internal/templater"
	"golang.org/x/exp/slog"
)

type Config struct {
	Provider string `yaml:"provider"`
	APIToken string `yaml:"api_token" env:"EMAIL_API_TOKEN"`
	Domain   string `yaml:"domain"`
}

type Template interface {
	Render(ctx context.Context, path string, params map[string]any) (*templater.Content, error)
}

type Idempotency interface {
	Get(ctx context.Context, key string) (bool, error)
	Save(ctx context.Context, key string) error
}

type Mailer interface {
	Send(ctx context.Context, target channel.Target, subject, body string) error
}

func New(config Config, template Template, idempotency Idempotency) *Email {
	var m Mailer
	provider := "noop"

	switch config.Provider {
	case "noop":
		fallthrough
	default:
		m = NewNoop()
	}

	slog.Debug("Email provider created", "provider", provider)

	return &Email{
		provider:    provider,
		mailer:      m,
		template:    template,
		idempotency: idempotency,
	}
}

type Email struct {
	provider    string
	mailer      Mailer
	template    Template
	idempotency Idempotency
}

func (m *Email) Send(ctx context.Context, templatePath string, target channel.Target, params map[string]any) error {
	ok, err := m.idempotency.Get(ctx, IdempotencyKey(target, templatePath))
	if err != nil {
		return err
	}

	if ok {
		slog.Info("Email skipped",
			"target", target.Address,
			"template", templatePath,
		)
		return nil
	}

	content, err := m.template.Render(ctx, templatePath, params)
	if err != nil {
		return err
	}

	if err := m.mailer.Send(ctx, target, content.Meta.Subject, content.Body); err != nil {
		return err
	}

	slog.Info("Email sent",
		"target", target.Address,
		"template", templatePath,
		"meta", content.Meta,
		"provider", m.provider,
	)

	k := IdempotencyKey(target, templatePath)
	if err := m.idempotency.Save(ctx, k); err != nil {
		return err
	}

	return nil
}

func (m *Email) Type() channel.Type {
	return channel.TypeEmail
}

func IdempotencyKey(target channel.Target, templatePath string) string {
	return fmt.Sprintf("%s:%s", target.Address, templatePath)
}
