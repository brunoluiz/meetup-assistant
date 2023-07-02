package repo

import (
	"context"
	"time"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
)

type EventsMock struct{}

func (m *EventsMock) GetActiveEvents(_ context.Context) ([]meetup_assistant.Event, error) {
	return []meetup_assistant.Event{
		{
			MeetupID: "meetup_id",
			Date:     time.Now().Add(24 * time.Hour * 7),
			Speakers: []meetup_assistant.Speaker{
				{
					Name:  "Bruno Luiz",
					Email: "bruno@bruno.com",
				},
			},
			Hosts: []meetup_assistant.Host{
				{
					Name:  "Bruno Luiz",
					Email: "bruno@bruno.com",
				},
			},
		},
	}, nil
}
