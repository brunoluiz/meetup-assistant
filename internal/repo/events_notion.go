package repo

import (
	"context"
	"fmt"
	"time"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
	"github.com/jomei/notionapi"
)

type EventsNotion struct {
	client *notionapi.Client
	db     notionapi.DatabaseID
	config meetup_assistant.NotionConfig
}

func NewNotion(token string, db string, config meetup_assistant.NotionConfig) (*EventsNotion, error) {
	m := &EventsNotion{
		client: notionapi.NewClient(notionapi.Token(token)),
		db:     notionapi.DatabaseID(db),
		config: config,
	}

	return m, m.Migrate(context.Background())
}

func (m *EventsNotion) Migrate(ctx context.Context) error {
	tables := m.config.Tables
	var err error

	_, err = m.client.Database.Update(ctx, tables.Submissions, &notionapi.DatabaseUpdateRequest{
		Title: []notionapi.RichText{
			{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: "🎤 Submissions"}},
		},
		Properties: notionapi.PropertyConfigs{
			"Name":        notionapi.TitlePropertyConfig{ID: "Name", Type: notionapi.PropertyConfigTypeTitle},
			"Email":       notionapi.RichTextPropertyConfig{ID: "Email", Type: notionapi.PropertyConfigTypeRichText},
			"Title":       notionapi.RichTextPropertyConfig{ID: "Title", Type: notionapi.PropertyConfigTypeRichText},
			"Description": notionapi.RichTextPropertyConfig{ID: "Description", Type: notionapi.PropertyConfigTypeRichText},
			"Bio":         notionapi.RichTextPropertyConfig{ID: "Bio", Type: notionapi.PropertyConfigTypeRichText},
			"Level":       notionapi.RichTextPropertyConfig{ID: "Level", Type: notionapi.PropertyConfigTypeRichText},
			"Format":      notionapi.RichTextPropertyConfig{ID: "Format", Type: notionapi.PropertyConfigTypeRichText},
			"Socials":     notionapi.RichTextPropertyConfig{ID: "Socials", Type: notionapi.PropertyConfigTypeRichText},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating submissions: %w", err)
	}

	_, err = m.client.Database.Update(ctx, tables.Venues, &notionapi.DatabaseUpdateRequest{
		Title: []notionapi.RichText{
			{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: "🏢 Venues"}},
		},
		Properties: notionapi.PropertyConfigs{
			"Name":    notionapi.TitlePropertyConfig{ID: "Name", Type: notionapi.PropertyConfigTypeTitle},
			"Address": notionapi.RichTextPropertyConfig{ID: "Address", Type: notionapi.PropertyConfigTypeRichText},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating venues: %w", err)
	}

	_, err = m.client.Database.Update(ctx, tables.Hosts, &notionapi.DatabaseUpdateRequest{
		Title: []notionapi.RichText{
			{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: "💁 Hosts"}},
		},
		Properties: notionapi.PropertyConfigs{
			"Name":  notionapi.TitlePropertyConfig{ID: "Name", Type: notionapi.PropertyConfigTypeTitle},
			"Email": notionapi.RichTextPropertyConfig{ID: "Email", Type: notionapi.PropertyConfigTypeRichText},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating hosts: %w", err)
	}

	_, err = m.client.Database.Update(ctx, tables.Events, &notionapi.DatabaseUpdateRequest{
		Title: []notionapi.RichText{
			{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: "📆 Events"}},
		},
		Properties: notionapi.PropertyConfigs{
			"Name":     notionapi.TitlePropertyConfig{ID: "Name", Type: notionapi.PropertyConfigTypeTitle},
			"MeetupID": notionapi.RichTextPropertyConfig{ID: "MeetupID", Type: notionapi.PropertyConfigTypeRichText},
			"Venue": notionapi.RelationPropertyConfig{
				Type: notionapi.PropertyConfigTypeRelation,
				Relation: notionapi.RelationConfig{
					DatabaseID:     tables.Venues,
					Type:           notionapi.RelationSingleProperty,
					SingleProperty: &notionapi.SingleProperty{},
				},
			},
			"Hosts": notionapi.RelationPropertyConfig{
				Type: notionapi.PropertyConfigTypeRelation,
				Relation: notionapi.RelationConfig{
					DatabaseID:     tables.Hosts,
					Type:           notionapi.RelationSingleProperty,
					SingleProperty: &notionapi.SingleProperty{},
				},
			},
			"Talks": notionapi.RelationPropertyConfig{
				Type: notionapi.PropertyConfigTypeRelation,
				Relation: notionapi.RelationConfig{
					DatabaseID:     tables.Submissions,
					Type:           notionapi.RelationSingleProperty,
					SingleProperty: &notionapi.SingleProperty{},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating events: %w", err)
	}

	return nil
}

func (m *EventsNotion) GetActiveEvents(ctx context.Context) ([]meetup_assistant.Event, error) {
	t := notionapi.Date(time.Now().Add(-7 * 24 * time.Hour))

	_, err := m.client.Database.Query(ctx, m.db, &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Date",
			Date: &notionapi.DateFilterCondition{
				OnOrAfter: &t,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return []meetup_assistant.Event{}, nil
}
