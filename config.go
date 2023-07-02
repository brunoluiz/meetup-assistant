package meetup_assistant

import (
	"fmt"
	"time"

	"github.com/brunoluiz/meetup-assistant/internal/tasker"
	"github.com/hako/durafmt"
	"github.com/jomei/notionapi"
)

type WhenType string

const (
	WhenTypeImmediate   WhenType = "immediate"
	WhenTypeBeforeEvent WhenType = "before"
	WhenTypeAfterEvent  WhenType = "after"
)

type NotionConfig struct {
	Tables NotionTables `json:"tables"`
}

type NotionTables struct {
	Events      notionapi.DatabaseID `json:"events"`
	Submissions notionapi.DatabaseID `json:"submissions"`
	Venues      notionapi.DatabaseID `json:"venues"`
	Hosts       notionapi.DatabaseID `json:"hosts"`
}

type DatabaseConfig struct {
	Notion NotionConfig `json:"notion" yaml:"notion"`
}

type Config struct {
	Comms    []CommJob      `yaml:"comms"`
	Database DatabaseConfig `yaml:"database"`
}

type CommJob struct {
	Audience string      `yaml:"audience"`
	Name     string      `yaml:"name"`
	Task     tasker.Task `yaml:"task"`
	Type     WhenType    `yaml:"type"`
	When     string      `yaml:"when"`
	Prev     string      `yaml:"prev"`
	Next     string      `yaml:"next"`
}

func (c *CommJob) Ready(prev string, now time.Time) (bool, error) {
	if prev != c.Prev {
		return false, nil
	}

	switch c.Type {
	case WhenTypeImmediate:
		return true, nil
	case WhenTypeBeforeEvent:
		t, err := durafmt.ParseString(c.When)
		if err != nil {
			return false, fmt.Errorf("invalid duration: %w", err)
		}

		beforeTime := now.Add(-1 * t.Duration())
		if now.After(beforeTime) {
			return true, nil
		}
	case WhenTypeAfterEvent:
		t, err := durafmt.ParseString(c.When)
		if err != nil {
			return false, fmt.Errorf("invalid duration: %w", err)
		}

		afterTime := now.Add(t.Duration())
		if now.After(afterTime) {
			return true, nil
		}
	}

	return false, nil
}
