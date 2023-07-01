package tasker

import (
	"context"

	"github.com/brunoluiz/meetup-assistant/internal/channel"
)

type Task struct {
	Channel  channel.Type   `yaml:"channel"`
	Template string         `yaml:"template"`
	Params   map[string]any `yaml:"params"`
}

type Channel interface {
	Send(ctx context.Context, templatePath string, params map[string]any) error
}

type ChannelFactory interface {
	Get(channelType channel.Type) (channel.Channel, error)
}

type Tasker struct {
	ChannelFactory ChannelFactory
}

func New(channelFactory ChannelFactory) *Tasker {
	return &Tasker{
		ChannelFactory: channelFactory,
	}
}

func (s *Tasker) Run(ctx context.Context, t Task, target channel.Target) error {
	c, err := s.ChannelFactory.Get(t.Channel)
	if err != nil {
		return err
	}

	return c.Send(ctx, t.Template, target, t.Params)
}
