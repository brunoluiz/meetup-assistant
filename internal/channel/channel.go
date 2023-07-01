package channel

import (
	"context"
	"fmt"
)

type Type string

const (
	TypeEmail    Type = "email"
	TypeTwitter  Type = "twitter"
	TypeLinkedIn Type = "linkedin"
)

type Target struct {
	Name    string
	Address string
	Meta    any
}

type Channel interface {
	Send(ctx context.Context, templatePath string, target Target, params map[string]any) error
	Type() Type
}

type Mailer interface {
	Send(ctx context.Context, subject string, to string, content string) error
}

type Factory struct {
	channel map[Type]Channel
}

func New(opts ...Channel) *Factory {
	f := &Factory{
		channel: map[Type]Channel{},
	}

	for _, opt := range opts {
		f.channel[opt.Type()] = opt
	}

	return f
}

func (f *Factory) Get(channelType Type) (Channel, error) {
	c, ok := f.channel[channelType]
	if !ok {
		return nil, fmt.Errorf("channel '%s' not supported", channelType)
	}

	return c, nil
}
