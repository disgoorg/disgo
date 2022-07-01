package rpc

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type ServerConfig struct {
	CDNHost     string `json:"cdn_host"`
	APIEndpoint string `json:"api_endpoint"`
	Environment string `json:"environment"`
}

type EventDataReady struct {
	V      int          `json:"v"`
	Config ServerConfig `json:"config"`
	User   discord.User `json:"user"`
}

func (EventDataReady) messageData() {}

var _ error = (*EventDataError)(nil)

type EventDataError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e EventDataError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (EventDataError) messageData() {}

type EventDataGuildStatus struct {
	Guild PartialGuild `json:"guild"`
}

func (EventDataGuildStatus) messageData() {}

type EventDataGuildCreate struct {
	ID   snowflake.ID `json:"id"`
	Name string       `json:"name"`
}

func (EventDataGuildCreate) messageData() {}

type EventDataChannelCreate struct {
	ID   snowflake.ID        `json:"id"`
	Name string              `json:"name"`
	Type discord.ChannelType `json:"type"`
}

func (EventDataChannelCreate) messageData() {}

type EventDataVoiceChannelSelect struct {
	ChannelID *snowflake.ID `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
}

func (EventDataVoiceChannelSelect) messageData() {}

type EventDataVoiceSettingsUpdate struct {
	VoiceSettings
}

func (EventDataVoiceSettingsUpdate) messageData() {}

type EventDataVoiceStateCreate struct {
	VoiceState
}

func (EventDataVoiceStateCreate) messageData() {}

type EventDataVoiceStateUpdate struct {
	VoiceState
}

func (EventDataVoiceStateUpdate) messageData() {}

type EventDataVoiceStateDelete struct {
	VoiceState
}

func (EventDataVoiceStateDelete) messageData() {}

type VoiceStateType string

const (
	VoiceStateTypeDisconnected      VoiceStateType = "DISCONNECTED"
	VoiceStateTypeAwaitingEndpoint  VoiceStateType = "AWAITING_ENDPOINT"
	VoiceStateTypeAuthenticating    VoiceStateType = "AUTHENTICATING"
	VoiceStateTypeConnecting        VoiceStateType = "CONNECTING"
	VoiceStateTypeConnected         VoiceStateType = "CONNECTED"
	VoiceStateTypeVoiceDisconnected VoiceStateType = "VOICE_DISCONNECTED"
	VoiceStateTypeVoiceConnecting   VoiceStateType = "VOICE_CONNECTING"
	VoiceStateTypeVoiceConnected    VoiceStateType = "VOICE_CONNECTED"
	VoiceStateTypeNoRoute           VoiceStateType = "NO_ROUTE"
	VoiceStateTypeICEChecking       VoiceStateType = "ICE_CHECKING"
)

type EventDataVoiceConnectionStatus struct {
	State       VoiceStateType `json:"state"`
	Hostname    string         `json:"hostname"`
	Pings       []float32      `json:"pings"`
	AveragePing float32        `json:"average_ping"`
	LastPing    float32        `json:"last_ping"`
}

func (EventDataVoiceConnectionStatus) messageData() {}

type EventDataMessageCreate struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Message   discord.Message `json:"message"`
}

func (EventDataMessageCreate) messageData() {}

type EventDataMessageUpdate struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Message   discord.Message `json:"message"`
}

func (EventDataMessageUpdate) messageData() {}

type EventDataMessageDelete struct {
	ChannelID snowflake.ID `json:"channel_id"`
	Message   struct {
		ID snowflake.ID `json:"id"`
	} `json:"message"`
}

func (EventDataMessageDelete) messageData() {}

type EventDataSpeakingStart struct {
	UserID snowflake.ID `json:"user_id"`
}

func (EventDataSpeakingStart) messageData() {}

type EventDataSpeakingStop struct {
	UserID snowflake.ID `json:"user_id"`
}

func (EventDataSpeakingStop) messageData() {}

type EventDataNotificationCreate struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Message   discord.Message `json:"message"`
	IconURL   string          `json:"icon_url"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
}

func (EventDataNotificationCreate) messageData() {}

type EventDataActivityJoin struct {
	Secret string `json:"secret"`
}

func (EventDataActivityJoin) messageData() {}

type EventDataActivitySpectate struct {
	Secret string `json:"secret"`
}

func (EventDataActivitySpectate) messageData() {}

type EventDataActivityJoinRequest struct {
	User discord.User `json:"user"`
}

func (EventDataActivityJoinRequest) messageData() {}
