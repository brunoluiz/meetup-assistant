package meetup_assistant

import (
	"fmt"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
)

func getTargets(audience string, e *Event) ([]channel.Target, error) {
	targets := []channel.Target{}

	switch audience {
	case "speakers":
		for _, s := range e.Speakers {
			targets = append(targets, channel.Target{Name: s.Name, Address: s.Email, Meta: s})
		}
		return targets, nil
	case "hosts":
		for _, h := range e.Hosts {
			targets = append(targets, channel.Target{Name: h.Name, Address: h.Email, Meta: h})
		}
		return targets, nil
	default:
		return nil, fmt.Errorf("event '%s' audience '%s' not supported", e.ID, audience)
	}
}
