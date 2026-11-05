package discord

import (
	"slices"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

// WebhookType (https: //discord.com/developers/docs/resources/webhook#webhook-object-webhook-types)
type WebhookType int

// All WebhookType(s)
const (
	WebhookTypeIncoming WebhookType = iota + 1
	WebhookTypeChannelFollower
	WebhookTypeApplication
)

// Webhook (https://discord.com/developers/docs/resources/webhook) is a way to post messages to Discord using the Discord API which do not require bot authentication or use.
type Webhook interface {
	json.Marshaler
	Type() WebhookType
	ID() snowflake.ID
	Name() string
	Avatar() *string
	AvatarURL(opts ...CDNOpt) *string
	CreatedAt() time.Time
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
		var v UnknownWebhook
		err = json.Unmarshal(data, &v)
		webhook = v
	}

	if err != nil {
		return err
	}

	w.Webhook = webhook
	return nil
}

var _ Webhook = (*IncomingWebhook)(nil)

type IncomingWebhook struct {
	id            snowflake.ID
	name          string
	avatar        *string
	ChannelID     snowflake.ID `json:"channel_id,omitempty"`
	GuildID       snowflake.ID `json:"guild_id,omitempty"`
	Token         string       `json:"token"`
	ApplicationID snowflake.ID `json:"application_id,omitempty"`
	User          *User        `json:"user,omitempty"`
	RawURL        *string      `json:"url,omitempty"`
}

func (w *IncomingWebhook) UnmarshalJSON(data []byte) error {
	type incomingWebhook IncomingWebhook
	var v struct {
		ID     snowflake.ID `json:"id"`
		Name   string       `json:"name"`
		Avatar *string      `json:"avatar"`
		incomingWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = IncomingWebhook(v.incomingWebhook)
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	return nil
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

func (w IncomingWebhook) ID() snowflake.ID {
	return w.id
}

func (w IncomingWebhook) Name() string {
	return w.name
}

func (w IncomingWebhook) Avatar() *string {
	return w.avatar
}

func (w IncomingWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if avatarURL := w.AvatarURL(opts...); avatarURL != nil {
		return *avatarURL
	}
	return w.DefaultAvatarURL(opts...)
}

func (w IncomingWebhook) AvatarURL(opts ...CDNOpt) *string {
	if w.Avatar() == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, w.ID(), *w.Avatar())
	return &url
}

func (w IncomingWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	return formatAssetURL(DefaultUserAvatar, opts, 0)
}

func (w IncomingWebhook) URL() string {
	return WebhookURL(w.ID(), w.Token)
}

func (w IncomingWebhook) CreatedAt() time.Time {
	return w.id.Time()
}

func (IncomingWebhook) webhook() {}

var _ Webhook = (*ChannelFollowerWebhook)(nil)

type ChannelFollowerWebhook struct {
	id            snowflake.ID
	name          string
	avatar        *string
	ChannelID     snowflake.ID          `json:"channel_id,omitempty"`
	GuildID       snowflake.ID          `json:"guild_id,omitempty"`
	SourceGuild   *WebhookSourceGuild   `json:"source_guild,omitempty"`
	SourceChannel *WebhookSourceChannel `json:"source_channel,omitempty"`
	User          *User                 `json:"user,omitempty"`
}

func (w *ChannelFollowerWebhook) UnmarshalJSON(data []byte) error {
	type channelFollowerWebhook ChannelFollowerWebhook
	var v struct {
		ID     snowflake.ID `json:"id"`
		Name   string       `json:"name"`
		Avatar *string      `json:"avatar"`
		channelFollowerWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = ChannelFollowerWebhook(v.channelFollowerWebhook)
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	return nil
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

func (w ChannelFollowerWebhook) ID() snowflake.ID {
	return w.id
}

func (w ChannelFollowerWebhook) Name() string {
	return w.name
}

func (w ChannelFollowerWebhook) Avatar() *string {
	return w.avatar
}

func (w ChannelFollowerWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if w.Avatar() == nil {
		return w.DefaultAvatarURL(opts...)
	}
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (w ChannelFollowerWebhook) AvatarURL(opts ...CDNOpt) *string {
	if w.Avatar() == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, w.ID(), *w.Avatar())
	return &url
}

func (w ChannelFollowerWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	return formatAssetURL(DefaultUserAvatar, opts, 0)
}

func (w ChannelFollowerWebhook) CreatedAt() time.Time {
	return w.id.Time()
}

func (ChannelFollowerWebhook) webhook() {}

var _ Webhook = (*ApplicationWebhook)(nil)

type ApplicationWebhook struct {
	id            snowflake.ID
	name          string
	avatar        *string
	ApplicationID snowflake.ID `json:"application_id,omitempty"`
	ChannelID     snowflake.ID `json:"channel_id,omitempty"`
	GuildID       snowflake.ID `json:"guild_id,omitempty"`
}

func (w *ApplicationWebhook) UnmarshalJSON(data []byte) error {
	type applicationWebhook ApplicationWebhook
	var v struct {
		ID     snowflake.ID `json:"id"`
		Name   string       `json:"name"`
		Avatar *string      `json:"avatar"`
		applicationWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = ApplicationWebhook(v.applicationWebhook)
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	return nil
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

func (w ApplicationWebhook) ID() snowflake.ID {
	return w.id
}

func (w ApplicationWebhook) Name() string {
	return w.name
}

func (w ApplicationWebhook) Avatar() *string {
	return w.avatar
}

func (w ApplicationWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if w.Avatar() == nil {
		return w.DefaultAvatarURL(opts...)
	}
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (w ApplicationWebhook) AvatarURL(opts ...CDNOpt) *string {
	if w.Avatar() == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, w.ID(), *w.Avatar())
	return &url
}

func (w ApplicationWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	return formatAssetURL(DefaultUserAvatar, opts, 0)
}

func (w ApplicationWebhook) CreatedAt() time.Time {
	return w.id.Time()
}

func (ApplicationWebhook) webhook() {}

var _ Webhook = (*UnknownWebhook)(nil)

// UnknownWebhook is a Webhook with a type the library does not recognize yet.
// Data holds the full payload, and MarshalJSON returns it when set.
type UnknownWebhook struct {
	id          snowflake.ID
	webhookType WebhookType
	name        string
	avatar      *string
	ChannelID   snowflake.ID    `json:"channel_id,omitempty"`
	GuildID     snowflake.ID    `json:"guild_id,omitempty"`
	Data        json.RawMessage `json:"-"`
}

func (w *UnknownWebhook) UnmarshalJSON(data []byte) error {
	type unknownWebhook UnknownWebhook
	var v struct {
		ID     snowflake.ID `json:"id"`
		Type   WebhookType  `json:"type"`
		Name   string       `json:"name"`
		Avatar *string      `json:"avatar"`
		unknownWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = UnknownWebhook(v.unknownWebhook)
	w.id = v.ID
	w.webhookType = v.Type
	w.name = v.Name
	w.avatar = v.Avatar
	w.Data = slices.Clone(data)
	return nil
}

func (w UnknownWebhook) MarshalJSON() ([]byte, error) {
	if len(w.Data) > 0 {
		return w.Data, nil
	}
	type unknownWebhook UnknownWebhook
	return json.Marshal(struct {
		ID     snowflake.ID `json:"id"`
		Type   WebhookType  `json:"type"`
		Name   string       `json:"name"`
		Avatar *string      `json:"avatar"`
		unknownWebhook
	}{
		ID:             w.id,
		Type:           w.webhookType,
		Name:           w.name,
		Avatar:         w.avatar,
		unknownWebhook: unknownWebhook(w),
	})
}

func (w UnknownWebhook) Type() WebhookType {
	return w.webhookType
}

func (w UnknownWebhook) ID() snowflake.ID {
	return w.id
}

func (w UnknownWebhook) Name() string {
	return w.name
}

func (w UnknownWebhook) Avatar() *string {
	return w.avatar
}

func (w UnknownWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return w.DefaultAvatarURL(opts...)
}

func (w UnknownWebhook) AvatarURL(opts ...CDNOpt) *string {
	if w.Avatar() == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, w.ID(), *w.Avatar())
	return &url
}

func (w UnknownWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	return formatAssetURL(DefaultUserAvatar, opts, 0)
}

func (w UnknownWebhook) CreatedAt() time.Time {
	return w.id.Time()
}

func (UnknownWebhook) webhook() {}

type WebhookSourceGuild struct {
	ID   snowflake.ID `json:"id"`
	Name string       `json:"name"`
	Icon *string      `json:"icon"`
}

type WebhookSourceChannel struct {
	ID   snowflake.ID `json:"id"`
	Name string       `json:"name"`
}

// WebhookCreate is used to create a Webhook
type WebhookCreate struct {
	Name   string `json:"name"`
	Avatar *Icon  `json:"avatar,omitempty"`
}

// WebhookUpdate is used to update a Webhook
type WebhookUpdate struct {
	Name      *string          `json:"name,omitempty"`
	Avatar    omit.Omit[*Icon] `json:"avatar,omitzero"`
	ChannelID *snowflake.ID    `json:"channel_id,omitempty"`
}

// WebhookUpdateWithToken is used to update a Webhook with the token
type WebhookUpdateWithToken struct {
	Name   *string          `json:"name,omitempty"`
	Avatar omit.Omit[*Icon] `json:"avatar,omitzero"`
}
