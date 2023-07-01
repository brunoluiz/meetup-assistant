package email

import (
	"context"

	"github.com/brunoluiz/meetup-assistant/internal/templater"
)

type Template interface {
	Render(ctx context.Context, path string, params map[string]any) (content *templater.Content, err error)
}

type Idempotency interface {
	Get(ctx context.Context, key string) (bool, error)
	Save(ctx context.Context, key string) error
}
