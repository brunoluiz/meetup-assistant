package main

import (
	"io"
	"os"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
	"github.com/brunoluiz/meetup-assistant/internal/channel"
	"github.com/brunoluiz/meetup-assistant/internal/channel/email"
	"github.com/brunoluiz/meetup-assistant/internal/source/git"
	"github.com/brunoluiz/meetup-assistant/internal/tasker"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	s, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	config := meetup_assistant.Config{}
	if err := yaml.Unmarshal(s, &config); err != nil {
		panic(err)
	}

	src, err := git.New(config.Template.Address)
	if err != nil {
		panic(err)
	}
	spew.Dump(src)

	channels := channel.New(
		email.NewNoop(nil, nil),
	)

	tasker := tasker.New(channels)
}
