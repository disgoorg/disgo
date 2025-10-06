package webhookevent

import (
	"fmt"
	"io"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type MessageType int

const (
	MessageTypePing MessageType = iota
	MessageTypeEvent
)

type Message struct {
	Version       int             `json:"version"`
	ApplicationID snowflake.ID    `json:"application_id"`
	Type          MessageType     `json:"type"`
	Event         MessageEvent    `json:"event,omitempty"`
	RawEvent      json.RawMessage `json:"-"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var v struct {
		Version       int             `json:"version"`
		ApplicationID snowflake.ID    `json:"application_id"`
		Type          MessageType     `json:"type"`
		Event         json.RawMessage `json:"event,omitempty"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		messageEvent MessageEvent
		err          error
	)

	switch v.Type {
	case MessageTypePing:
		messageEvent = EventPing{}
	case MessageTypeEvent:
		var d Event
		err = json.Unmarshal(v.Event, &d)
		messageEvent = d
	default:
		var d EventUnknown
		err = json.Unmarshal(v.Event, &d)
		messageEvent = d
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal message data: %s: %w", string(data), err)
	}
	m.Version = v.Version
	m.ApplicationID = v.ApplicationID
	m.Type = v.Type
	m.Event = messageEvent
	m.RawEvent = v.Event
	return nil
}

type MessageEvent interface {
	messageEvent()
}

type EventPing struct{}

func (EventPing) messageEvent() {}

type EventUnknown json.RawMessage

func (EventUnknown) messageEvent() {}

type EventType string

const (
	EventTypeApplicationAuthorized   EventType = "APPLICATION_AUTHORIZED"
	EventTypeApplicationDeauthorized EventType = "APPLICATION_DEAUTHORIZED"
	EventTypeEntitlementCreate       EventType = "ENTITLEMENT_CREATE"
	EventTypeQuestUserEnrollment     EventType = "QUEST_USER_ENROLLMENT"
)

type Event struct {
	Type      EventType       `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      EventData       `json:"data"`
	RawData   json.RawMessage `json:"-"`
}

func (Event) messageEvent() {}

func (m *Event) UnmarshalJSON(data []byte) error {
	var v struct {
		Type      EventType       `json:"type"`
		Timestamp time.Time       `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		eventData EventData
		err       error
	)

	switch v.Type {
	case EventTypeApplicationAuthorized:
		var d EventDataApplicationAuthorized
		err = json.Unmarshal(v.Data, &d)
		eventData = d
	case EventTypeApplicationDeauthorized:
		var d EventDataApplicationDeauthorized
		err = json.Unmarshal(v.Data, &d)
		eventData = d
	case EventTypeEntitlementCreate:
		var d EventDataEntitlementCreate
		err = json.Unmarshal(v.Data, &d)
		eventData = d
	case EventTypeQuestUserEnrollment:
		var d EventDataQuestUserEnrollment
		err = json.Unmarshal(v.Data, &d)
		eventData = d
	default:
		var d EventDataUnknown
		err = json.Unmarshal(v.Data, &d)
		eventData = d
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal message data: %s: %w", string(data), err)
	}
	m.Type = v.Type
	m.Timestamp = v.Timestamp
	m.Data = eventData
	m.RawData = v.Data
	return nil
}

type EventData interface {
	eventData()
}

type EventDataApplicationAuthorized struct {
	IntegrationType discord.IntegrationType `json:"integration_type"`
	User            discord.User            `json:"user"`
	Scopes          []discord.OAuth2Scope   `json:"scopes"`
	Guild           discord.Guild           `json:"guild"`
}

func (EventDataApplicationAuthorized) eventData() {}

type EventDataApplicationDeauthorized struct {
	User discord.User `json:"user"`
}

func (EventDataApplicationDeauthorized) eventData() {}

type EventDataEntitlementCreate struct {
	discord.Entitlement
}

func (EventDataEntitlementCreate) eventData() {}

type EventDataQuestUserEnrollment struct {
	// This event cannot be received by apps at this time. It's documented because it appears on the Webhooks settings page.
}

func (EventDataQuestUserEnrollment) eventData() {}

type EventDataUnknown json.RawMessage

func (EventDataUnknown) eventData() {}

type EventDataRaw struct {
	EventType EventType
	Payload   io.Reader
}

func (EventDataRaw) eventData() {}
