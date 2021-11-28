package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

// WebhookType (https: //discord.com/developers/docs/resources/webhook#webhook-object-webhook-types)
type WebhookType int

// All WebhookType(s)
//goland:noinspection GoUnusedConst
const (
	WebhookTypeIncoming WebhookType = iota + 1
	WebhookTypeChannelFollower
	WebhookTypeApplication
)

// Webhook (https://discord.com/developers/docs/resources/webhook) is a way to post messages to Discord using the Discord API which do not require bot authentication or use.
type Webhook interface {
	json.Marshaler
	Type() WebhookType
	ID() Snowflake
	webhook()
}

type UnmarshalWebhook struct {
	Webhook
}

func (w *UnmarshalWebhook) UnmarshalJSON(data []byte) error {
	var wType struct {
		Type WebhookType `json:"type"`
	}

	if err := json.Unmarshal(data, &wType); err != nil {
		return err
	}

	var (
		webhook Webhook
		err     error
	)

	switch wType.Type {
	case WebhookTypeIncoming:
		var v IncomingWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	case WebhookTypeChannelFollower:
		var v ChannelFollowerWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	case WebhookTypeApplication:
		var v ApplicationWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	default:
		err = fmt.Errorf("unkown webhook with type %d received", wType.Type)
	}

	if err != nil {
		return err
	}

	w.Webhook = webhook
	return nil
}

var _ Webhook = (*IncomingWebhook)(nil)

type IncomingWebhook struct {
	WebhookID     Snowflake  `json:"id"`
	Name          string     `json:"name"`
	Avatar        *string    `json:"avatar"`
	ChannelID     Snowflake  `json:"channel_id"`
	GuildID       Snowflake  `json:"guild_id"`
	Token         string     `json:"token"`
	ApplicationID *Snowflake `json:"application_id"`
	User          User       `json:"user"`
}

func (w IncomingWebhook) MarshalJSON() ([]byte, error) {
	type incomingWebhook IncomingWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		incomingWebhook
	}{
		Type:            w.Type(),
		incomingWebhook: incomingWebhook(w),
	})
}

func (IncomingWebhook) Type() WebhookType {
	return WebhookTypeIncoming
}

func (IncomingWebhook) webhook() {}

func (w IncomingWebhook) ID() Snowflake {
	return w.WebhookID
}

var _ Webhook = (*ChannelFollowerWebhook)(nil)

type ChannelFollowerWebhook struct {
	WebhookID     Snowflake            `json:"id"`
	Name          string               `json:"name"`
	Avatar        *string              `json:"avatar"`
	ChannelID     Snowflake            `json:"channel_id"`
	GuildID       Snowflake            `json:"guild_id"`
	SourceGuild   WebhookSourceGuild   `json:"source_guild"`
	SourceChannel WebhookSourceChannel `json:"source_channel"`
	User          User                 `json:"user"`
}

func (w ChannelFollowerWebhook) MarshalJSON() ([]byte, error) {
	type channelFollowerWebhook ChannelFollowerWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		channelFollowerWebhook
	}{
		Type:                   w.Type(),
		channelFollowerWebhook: channelFollowerWebhook(w),
	})
}

func (ChannelFollowerWebhook) Type() WebhookType {
	return WebhookTypeChannelFollower
}

func (ChannelFollowerWebhook) webhook() {}

func (w ChannelFollowerWebhook) ID() Snowflake {
	return w.WebhookID
}

var _ Webhook = (*ApplicationWebhook)(nil)

type ApplicationWebhook struct {
	WebhookID     Snowflake `json:"id"`
	Name          string    `json:"name"`
	Avatar        *string   `json:"avatar"`
	ApplicationID Snowflake `json:"application_id"`
}

func (w ApplicationWebhook) MarshalJSON() ([]byte, error) {
	type applicationWebhook ApplicationWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		applicationWebhook
	}{
		Type:               w.Type(),
		applicationWebhook: applicationWebhook(w),
	})
}

func (ApplicationWebhook) Type() WebhookType {
	return WebhookTypeApplication
}

func (ApplicationWebhook) webhook() {}

func (w ApplicationWebhook) ID() Snowflake {
	return w.WebhookID
}

type WebhookSourceGuild struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
	Icon *string   `json:"icon"`
}

type WebhookSourceChannel struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
}

// WebhookCreate is used to create a Webhook
type WebhookCreate struct {
	Name   string `json:"name"`
	Avatar *Icon  `json:"avatar,omitempty"`
}

// WebhookUpdate is used to update a Webhook
type WebhookUpdate struct {
	Name      *string    `json:"name,omitempty"`
	Avatar    *NullIcon  `json:"avatar,omitempty"`
	ChannelID *Snowflake `json:"channel_id"`
}

// WebhookUpdateWithToken is used to update a Webhook with the token
type WebhookUpdateWithToken struct {
	Name   *string   `json:"name,omitempty"`
	Avatar *NullIcon `json:"avatar,omitempty"`
}
