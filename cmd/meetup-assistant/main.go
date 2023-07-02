package main

import (
	"context"
	"os"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/brunoluiz/meetup-assistant/internal/channel/email"
	"github.com/brunoluiz/meetup-assistant/internal/repo"
	"github.com/brunoluiz/meetup-assistant/internal/source/git"
	"github.com/brunoluiz/meetup-assistant/internal/storage"
	"github.com/brunoluiz/meetup-assistant/internal/tasker"
	"github.com/brunoluiz/meetup-assistant/internal/templater"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

func run() error {
	ctx := context.Background()

	s, err := os.ReadFile("sample.yaml")
	if err != nil {
		return err
	}

	config := meetup_assistant.Config{}
	if err := yaml.Unmarshal(s, &config); err != nil {
		return err
	}

	src, err := git.New(config.Template.Address)
	if err != nil {
		return err
	}

	idempotency, err := storage.NewIdempotencyJSON("/tmp/idempotency.json")
	if err != nil {
		return err
	}

	state, err := storage.NewStateJSON("/tmp/state.json")
	if err != nil {
		return err
	}

	channels := channel.New(
		email.New(
			config.Email,
			templater.NewMarkdownHTML(src),
			idempotency,
		),
	)

	tasker := tasker.New(channels)
	r, err := repo.New(os.Getenv("DB_DSN"), config.Database)
	if err != nil {
		return err
	}

	return meetup_assistant.New(r, state, tasker, config.Comms).Run(ctx)
}

func main() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
