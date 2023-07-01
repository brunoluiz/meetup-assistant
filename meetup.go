package meetup_assistant

import (
	"context"
	"errors"
	"time"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/brunoluiz/meetup-assistant/internal/tasker"
)

type Speaker struct {
	Name  string
	Email string
}

type Host struct {
	Name  string
	Email string
}

type Event struct {
	ID       string
	Date     time.Time
	Speakers []Speaker
	Hosts    []Host
}

type Repository interface {
	GetActiveEvents(ctx context.Context) ([]*Event, error)
}

type StateStorage interface {
	Get(ctx context.Context, audience, key string) (string, error)
	Save(ctx context.Context, audience, key, state string) error
}

type Tasker interface {
	Run(ctx context.Context, t tasker.Task, target channel.Target) error
}

type Meetup struct {
	repo   Repository
	state  StateStorage
	tasker Tasker
	jobs   []CommJob
}

func (s *Meetup) Run(ctx context.Context) error {
	events, err := s.repo.GetActiveEvents(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	var errs error

	for _, event := range events {
		for _, config := range s.jobs {
			targets, err := getTargets(config.Audience, event)
			if err != nil {
				return err
			}

			for _, target := range targets {
				prevState, err := s.state.Get(ctx, config.Audience, target.Address)
				if err != nil {
					errs = errors.Join(errs, err)
				}

				if ok, err := config.Ready(prevState, now); !ok {
					if err != nil {
						err = errors.Join(errs, err)
					}
					continue
				}

				if err := s.tasker.Run(ctx, config.Task, target); err != nil {
					errs = errors.Join(errs, err)
				}
			}
		}
	}

	return err
}
