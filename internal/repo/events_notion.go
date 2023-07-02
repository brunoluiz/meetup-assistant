package repo

import (
	"context"
	"fmt"
	"time"

	meetup_assistant "github.com/brunoluiz/meetup-assistant"
	"github.com/davecgh/go-spew/spew"
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

	return m, nil
	// return m, m.Migrate(context.Background())
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
			"Date":     notionapi.DatePropertyConfig{ID: "Date", Type: notionapi.PropertyConfigTypeDate},
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

	res, err := m.client.Database.Query(ctx, m.config.Tables.Events, &notionapi.DatabaseQueryRequest{
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

	out := []meetup_assistant.Event{}

	for _, e := range res.Results {
		event := meetup_assistant.Event{}
		event.MeetupID = e.Properties["MeetupID"].(*notionapi.RichTextProperty).RichText[0].PlainText
		event.Name = e.Properties["Name"].(*notionapi.TitleProperty).Title[0].PlainText
		event.Date = time.Time(*e.Properties["Date"].(*notionapi.DateProperty).Date.Start)
		// spew.Dump(e.Properties)

		event.Hosts, err = m.getHosts(ctx, e.Properties["Hosts"].(*notionapi.RelationProperty).Relation[0].ID)
		if err != nil {
			return nil, err
		}

		event.Venue, err = m.getVenue(ctx, e.Properties["Venue"].(*notionapi.RelationProperty).Relation[0].ID)
		if err != nil {
			return nil, err
		}

		event.Speakers, err = m.getSubmissions(ctx, e.Properties["Talks"].(*notionapi.RelationProperty).Relation[0].ID)
		if err != nil {
			return nil, err
		}

		out = append(out, event)
	}
	spew.Dump(out)

	return out, nil
}

func (m *EventsNotion) getByName(ctx context.Context, table notionapi.DatabaseID, ids ...notionapi.PageID) ([]notionapi.Properties, error) {
	out := []notionapi.Properties{}
	for _, id := range ids {
		p, err := m.client.Page.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		out = append(out, p.Properties)
	}

	return out, nil
}

func (m *EventsNotion) getHosts(ctx context.Context, ids ...notionapi.PageID) ([]meetup_assistant.Host, error) {
	res, err := m.getByName(ctx, m.config.Tables.Hosts, ids...)
	if err != nil {
		return nil, err
	}

	var out []meetup_assistant.Host
	for _, p := range res {
		obj := meetup_assistant.Host{}
		obj.Name = p["Name"].(*notionapi.TitleProperty).Title[0].PlainText
		obj.Email = p["Email"].(*notionapi.RichTextProperty).RichText[0].PlainText

		out = append(out, obj)
	}

	return out, nil
}

func (m *EventsNotion) getSubmissions(ctx context.Context, ids ...notionapi.PageID) ([]meetup_assistant.Speaker, error) {
	res, err := m.getByName(ctx, m.config.Tables.Submissions, ids...)
	if err != nil {
		return nil, err
	}

	var out []meetup_assistant.Speaker
	for _, p := range res {
		obj := meetup_assistant.Speaker{}
		obj.Name = p["Name"].(*notionapi.TitleProperty).Title[0].PlainText
		obj.Email = p["Email"].(*notionapi.RichTextProperty).RichText[0].PlainText

		out = append(out, obj)
	}

	return out, nil
}

func (m *EventsNotion) getVenue(ctx context.Context, ids ...notionapi.PageID) (meetup_assistant.Venue, error) {
	res, err := m.getByName(ctx, m.config.Tables.Submissions, ids...)
	if err != nil {
		return meetup_assistant.Venue{}, err
	}

	for _, p := range res {
		obj := meetup_assistant.Venue{}
		obj.Name = p["Name"].(*notionapi.TitleProperty).Title[0].PlainText
		obj.Address = p["Address"].(*notionapi.RichTextProperty).RichText[0].PlainText
		return obj, nil
	}

	return meetup_assistant.Venue{}, fmt.Errorf("no venue is registered")
}
