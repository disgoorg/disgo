package rpc

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type CmdArgsAuthorize struct {
	ClientID snowflake.ID          `json:"client_id"`
	Scopes   []discord.OAuth2Scope `json:"scopes"`
	RPCToken string                `json:"rpc_token,omitempty"`
	Username string                `json:"username,omitempty"`
}

func (CmdArgsAuthorize) cmdArgs() {}

type CmdRsAuthorize struct {
	Code string `json:"code"`
}

func (CmdRsAuthorize) messageData() {}

type CmdArgsAuthenticate struct {
	AccessToken string `json:"access_token"`
}

func (CmdArgsAuthenticate) cmdArgs() {}

type CmdRsAuthenticate struct {
	User        discord.User              `json:"user"`
	Scopes      []discord.OAuth2Scope     `json:"scopes"`
	Expires     time.Time                 `json:"expires"`
	Application discord.OAuth2Application `json:"application"`
}

func (CmdRsAuthenticate) messageData() {}

type CmdArgsGetGuild struct {
	GuildID snowflake.ID `json:"guild_id"`
	Timeout int          `json:"timeout"`
}

func (CmdArgsGetGuild) cmdArgs() {}

type PartialGuild struct {
	ID      snowflake.ID `json:"id"`
	Name    string       `json:"name"`
	IconURL *string      `json:"icon_url,omitempty"`
}

type CmdRsGetGuild struct {
	PartialGuild
}

func (CmdRsGetGuild) messageData() {}

type CmdRsGetGuilds struct {
	Guilds []PartialGuild `json:"guilds"`
}

func (CmdRsGetGuilds) messageData() {}

type CmdArgsGetChannel struct {
	ChannelID snowflake.ID `json:"channel_id"`
}

func (CmdArgsGetChannel) cmdArgs() {}

type PartialChannel struct {
	ID          snowflake.ID        `json:"id"`
	GuildID     *snowflake.ID       `json:"guild_id,omitempty"`
	Name        string              `json:"name"`
	Type        discord.ChannelType `json:"type"`
	Topic       *string             `json:"topic,omitempty"`
	Bitrate     int                 `json:"bitrate,omitempty"`
	UserLimit   int                 `json:"user_limit,omitempty"`
	Position    int                 `json:"position,omitempty"`
	VoiceStates []VoiceState        `json:"voice_states,omitempty"`
	Messages    []ChannelMessage    `json:"messages,omitempty"`
}

type VoiceState struct {
	discord.VoiceState
	Volume int `json:"volume"`
	Pan    Pan `json:"pan"`
}

type Pan struct {
	Left  float32 `json:"left"`
	Right float32 `json:"right"`
}

type CmdRsGetChannel struct {
	PartialChannel
}

func (CmdRsGetChannel) messageData() {}

type CmdArgsGetChannels struct {
	GuildID snowflake.ID `json:"guild_id"`
}

func (CmdArgsGetChannels) cmdArgs() {}

type CmdRsGetChannels struct {
	Channels []PartialChannel `json:"channels"`
}

func (CmdRsGetChannels) messageData() {}

type CmdArgsSetUserVoiceSettings struct {
	UserID snowflake.ID `json:"user_id"`
	Pan    *Pan         `json:"pan,omitempty"`
	Volume *int         `json:"volume,omitempty"`
	Mute   *bool        `json:"mute,omitempty"`
}

func (CmdArgsSetUserVoiceSettings) cmdArgs() {}

type CmdRsSetUserVoiceSettings struct {
	UserID snowflake.ID `json:"user_id"`
	Pan    Pan          `json:"pan"`
	Volume int          `json:"volume"`
	Mute   bool         `json:"mute"`
}

func (CmdRsSetUserVoiceSettings) messageData() {}

type CmdArgsSelectVoiceChannel struct {
	ChannelID snowflake.ID `json:"channel_id"`
	Timeout   int          `json:"timeout"`
	Force     bool         `json:"force"`
}

func (CmdArgsSelectVoiceChannel) cmdArgs() {}

type CmdRsSelectVoiceChannel struct {
	PartialChannel
}

func (CmdRsSelectVoiceChannel) messageData() {}

type CmdRsGetSelectedVoiceChannel struct {
	*PartialChannel
}

func (CmdRsGetSelectedVoiceChannel) messageData() {}

type CmdArgsSelectTextChannel struct {
	ChannelID *snowflake.ID `json:"channel_id"`
	Timeout   int           `json:"timeout"`
}

func (CmdArgsSelectTextChannel) cmdArgs() {}

type CmdRsSelectTextChannel struct {
	*PartialChannel
}

func (CmdRsSelectTextChannel) messageData() {}

type CmdRsGetVoiceSettings struct {
	VoiceSettings
}

type VoiceSettings struct {
	Input                VoiceSettingsIO   `json:"input"`
	Output               VoiceSettingsIO   `json:"output"`
	Mode                 VoiceSettingsMode `json:"mode"`
	AutomaticGainControl bool              `json:"automatic_gain_control"`
	EchoCancellation     bool              `json:"echo_cancellation"`
	NoiseSuppression     bool              `json:"noise_suppression"`
	QOS                  bool              `json:"qos"`
	SilenceWarning       bool              `json:"silence_warning"`
	Deaf                 bool              `json:"deaf"`
	Mute                 bool              `json:"mute"`
}

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VoiceSettingsIO struct {
	DeviceID         string   `json:"device_id"`
	Volume           int      `json:"volume"`
	AvailableDevices []Device `json:"available_devices"`
}

type VoiceSettingsModeType string

const (
	VoiceSettingsModeTypePushToTalk VoiceSettingsModeType = "PUSH_TO_TALK"
	VoiceSettingsModeTypeActivity   VoiceSettingsModeType = "ACTIVITY"
)

type VoiceSettingsMode struct {
	Type          VoiceSettingsModeType `json:"type"`
	AutoThreshold bool                  `json:"auto_threshold"`
	Threshold     int                   `json:"threshold"`
	Shortcut      ShortcutKeyCombo      `json:"shortcut"`
	Delay         float32               `json:"delay"`
}

type ShortcutKeyComboType int

const (
	ShortcutKeyComboTypeKeyboardKey ShortcutKeyComboType = iota
	ShortcutKeyComboTypeMouseButton
	ShortcutKeyComboTypeModifierKey
	ShortcutKeyComboTypeGamepadButton
)

type ShortcutKeyCombo struct {
	Type ShortcutKeyComboType `json:"type"`
	Code int                  `json:"code"`
	Name string               `json:"name"`
}

func (CmdRsGetVoiceSettings) messageData() {}

type CmdArgsSetVoiceSettings struct {
	Input                *VoiceSettings     `json:"input"`
	Output               *VoiceSettings     `json:"output"`
	Mode                 *VoiceSettingsMode `json:"mode"`
	AutomaticGainControl *bool              `json:"automatic_gain_control"`
	EchoCancellation     *bool              `json:"echo_cancellation"`
	NoiseSuppression     *bool              `json:"noise_suppression"`
	QOS                  *bool              `json:"qos"`
	SilenceWarning       *bool              `json:"silence_warning"`
	Deaf                 *bool              `json:"deaf"`
	Mute                 *bool              `json:"mute"`
}

func (CmdArgsSetVoiceSettings) cmdArgs() {}

type CmdRsSetVoiceSettings struct {
	VoiceSettings
}

func (CmdRsSetVoiceSettings) messageData() {}

type DeviceType string

const (
	DeviceTypeAudioInput  DeviceType = "audioinput"
	DeviceTypeAudioOutput DeviceType = "audiooutput"
	DeviceTypeVideoInput  DeviceType = "videoinput"
)

type DeviceVendor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type DeviceModel struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CertifiedDevice struct {
	Type                 DeviceType   `json:"type"`
	ID                   string       `json:"id"`
	Vendor               DeviceVendor `json:"vendor"`
	Model                DeviceModel  `json:"model"`
	Related              []string     `json:"related"`
	EchoCancellation     bool         `json:"echo_cancellation"`
	NoiseSuppression     bool         `json:"noise_suppression"`
	AutomaticGainControl bool         `json:"automatic_gain_control"`
	HardwareMute         bool         `json:"hardware_mute"`
}

type CmdArgsSetCertifiedDevices struct {
	Devices []CertifiedDevice `json:"devices"`
}

func (CmdArgsSetCertifiedDevices) cmdArgs() {}

type CmdArgsSendActivityJoinInvite struct {
	UserID snowflake.ID `json:"user_id"`
}

func (CmdArgsSendActivityJoinInvite) cmdArgs() {}

type CmdArgsCloseActivityRequest struct {
	UserID snowflake.ID `json:"user_id"`
}

func (CmdArgsCloseActivityRequest) cmdArgs() {}

type CmdArgsSetActivity struct {
	PID      int              `json:"pid"`
	Activity discord.Activity `json:"activity"`
}

func (CmdArgsSetActivity) cmdArgs() {}

type CmdRsSetActivity struct {
	discord.Activity
}

func (CmdRsSetActivity) messageData() {}

type CmdArgsSubscribe interface {
	CmdArgs
	cmdArgsSubscribe()
}

type CmdArgsSubscribeMessage struct {
	ChannelID snowflake.ID `json:"channel_id"`
}

func (CmdArgsSubscribeMessage) cmdArgs()          {}
func (CmdArgsSubscribeMessage) cmdArgsSubscribe() {}

type CmdArgsSubscribeGuild struct {
	GuildID snowflake.ID `json:"guild_id"`
}

func (CmdArgsSubscribeGuild) cmdArgs()          {}
func (CmdArgsSubscribeGuild) cmdArgsSubscribe() {}

type CmdArgsSubscribeSpeaking struct {
	ChannelID snowflake.ID `json:"channel_id"`
}

func (CmdArgsSubscribeSpeaking) cmdArgs()          {}
func (CmdArgsSubscribeSpeaking) cmdArgsSubscribe() {}

type CmdRsSubscribe struct {
	Evt string `json:"evt"`
}

func (CmdRsSubscribe) messageData() {}

type CmdRsUnsubscribe struct {
	Evt string `json:"evt"`
}

func (CmdRsUnsubscribe) messageData() {}
