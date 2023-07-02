package repo

import (
	"context"
	"fmt"
	"net/url"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
)

type Repository interface {
	GetActiveEvents(ctx context.Context) ([]meetup_assistant.Event, error)
}

func New(dsn string, config meetup_assistant.DatabaseConfig) (Repository, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "notion":
		token := u.Query().Get("token")
		db := u.Hostname()
		return NewNotion(token, db, config.Notion)
	case "mock":
		return &EventsMock{}, nil
	default:
		return nil, fmt.Errorf("unknown scheme: '%s'", u.Scheme)
	}
}
