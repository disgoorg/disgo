package events

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
)

// NewEvent constructs a new event with the provided Client instance
func NewEvent(client *bot.Client) *Event {
	return &Event{client: client}
}

// event the base event structure
type Event struct {
	client *bot.Client
}

// Client returns the bot.Client instance that dispatched the event
func (e *Event) Client() *bot.Client {
	return e.client
}

func NewGatewayEvent(sequenceNumber int, shardID int) *GatewayEvent {
	return &GatewayEvent{
		sequenceNumber: sequenceNumber,
		shardID:        shardID,
	}
}

type GatewayEvent struct {
	sequenceNumber int
	shardID        int
}

// SequenceNumber returns the sequence number of the gateway event
func (e *GatewayEvent) SequenceNumber() int {
	return e.sequenceNumber
}

// ShardID returns the shard ID the event was dispatched from
func (e *GatewayEvent) ShardID() int {
	return e.shardID
}

// NewWebhookEvent constructs a new WebhookEvent with the provided Client instance and timestamp
func NewWebhookEvent(version int, applicationID snowflake.ID, timestamp time.Time) *WebhookEvent {
	return &WebhookEvent{
		version:       version,
		applicationID: applicationID,
		timestamp:     timestamp,
	}
}

// WebhookEvent represents an event received from a webhook
type WebhookEvent struct {
	version       int
	applicationID snowflake.ID
	timestamp     time.Time
}

// Version returns the version of the webhook event
func (e *WebhookEvent) Version() int {
	return e.version
}

// ApplicationID returns the application ID associated with the webhook event
func (e *WebhookEvent) ApplicationID() snowflake.ID {
	return e.applicationID
}

// Timestamp returns the timestamp of the webhook event
func (e *WebhookEvent) Timestamp() time.Time {
	return e.timestamp
}
