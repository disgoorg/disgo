package rpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type AvatarDecorationData struct {
	Asset string `json:"asset"`
	SkuID string `json:"sku_id"`
}

// User is a struct for interacting with discord's users
type User struct {
	discord.User
}

func (u *User) UnmarshalJSON(data []byte) error {
	type user User
	var v struct {
		user
		Flags                discord.UserFlags    `json:"flags"`
		AvatarDecorationData AvatarDecorationData `json:"avatar_decoration_data"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = User(v.user)
	u.PublicFlags = v.Flags
	u.AvatarDecoration = &v.AvatarDecorationData.Asset
	return nil
}

// Attachment is used for files sent in a Message
type Attachment struct {
	discord.Attachment
	Spoiler bool `json:"spoiler,omitempty"`
}

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	discord.Message
	Attachments []Attachment   `json:"attachments"`
	Author      User           `json:"author"`
	Blocked     bool           `json:"blocked"`
	Bot         bool           `json:"bot"`
	Embeds      []Embed        `json:"embeds"`
	Mentions    []snowflake.ID `json:"mentions"`
	Nick        string         `json:"nick"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v struct {
		message
		Timestamp time.Time `json:"timestamp"` // discord.Message.CreatedAt
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = Message(v.message)
	m.CreatedAt = v.Timestamp
	m.Attachments = v.Attachments
	m.Blocked = v.Blocked
	m.Bot = v.Bot
	m.Embeds = v.Embeds
	m.Mentions = v.Mentions
	m.Nick = v.Nick
	return nil
}

type message struct {
	Cmd  Cmd     `json:"cmd"`
	Args CmdArgs `json:"args,omitempty"`

	Event EventType   `json:"evt,omitempty"`
	Data  MessageData `json:"data,omitempty"`

	Nonce string `json:"nonce,omitempty"`
}

func (m *message) UnmarshalJSON(data []byte) error {
	type msg message
	var v struct {
		Data json.RawMessage `json:"data"`
		msg
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	m.Cmd = v.Cmd
	m.Event = v.Event
	m.Nonce = v.Nonce

	var (
		messageData MessageData
		err         error
	)

	if v.Event == EventError {
		var d EventDataError
		err = json.Unmarshal(v.Data, &d)
		messageData = d
	} else {
		switch v.Cmd {
		case CmdDispatch:
			switch v.Event {
			case EventReady:
				var d EventDataReady
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventGuildStatus:
				var d EventDataGuildStatus
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventGuildCreate:
				var d EventDataGuildCreate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventChannelCreate:
				var d EventDataChannelCreate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceChannelSelect:
				var d EventDataVoiceChannelSelect
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceSettingsUpdate:
				var d EventDataVoiceSettingsUpdate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceStateCreate:
				var d EventDataVoiceStateCreate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceStateUpdate:
				var d EventDataVoiceStateUpdate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceStateDelete:
				var d EventDataVoiceStateDelete
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventVoiceConnectionStatus:
				var d EventDataVoiceConnectionStatus
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventMessageCreate:
				var d EventDataMessageCreate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventMessageUpdate:
				var d EventDataMessageUpdate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventMessageDelete:
				var d EventDataMessageDelete
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventSpeakingStart:
				var d EventDataSpeakingStart
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventSpeakingStop:
				var d EventDataSpeakingStop
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventNotificationCreate:
				var d EventDataNotificationCreate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventActivityJoin:
				var d EventDataActivityJoin
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventActivitySpectate:
				var d EventDataActivitySpectate
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			case EventActivityJoinRequest:
				var d EventDataActivityJoinRequest
				err = json.Unmarshal(v.Data, &d)
				messageData = d

			default:
				err = fmt.Errorf("unknown event: %s", v.Event)
			}

		case CmdAuthorize:
			var d CmdRsAuthorize
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdAuthenticate:
			var d CmdRsAuthenticate
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetGuild:
			var d CmdRsGetGuild
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetGuilds:
			var d CmdRsGetGuilds
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetChannel:
			var d CmdRsGetChannel
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetChannels:
			var d CmdRsGetChannels
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSetUserVoiceSettings:
			var d CmdRsSetUserVoiceSettings
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSelectVoiceChannel:
			var d CmdRsSelectVoiceChannel
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetSelectedVoiceChannel:
			var d CmdRsGetSelectedVoiceChannel
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSelectTextChannel:
			var d CmdRsSelectTextChannel
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdGetVoiceSettings:
			var d CmdRsGetVoiceSettings
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSetVoiceSettings:
			var d CmdRsSetVoiceSettings
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSetCertifiedDevices:
		// no response data

		case CmdSendActivityJoinInvite:
		// no response data

		case CmdCloseActivityRequest:
		// no response data

		case CmdSubscribe:
			var d CmdRsSubscribe
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdUnsubscribe:
			var d CmdRsUnsubscribe
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		case CmdSetActivity:
			var d CmdRsSetActivity
			err = json.Unmarshal(v.Data, &d)
			messageData = d

		default:
			err = fmt.Errorf("unknown cmd: %s", v.Cmd)
		}
	}

	if err != nil {
		return err
	}

	m.Data = messageData
	return nil
}

type Cmd string

const (
	CmdDispatch                Cmd = "DISPATCH"
	CmdAuthorize               Cmd = "AUTHORIZE"
	CmdAuthenticate            Cmd = "AUTHENTICATE"
	CmdGetGuild                Cmd = "GET_GUILD"
	CmdGetGuilds               Cmd = "GET_GUILDS"
	CmdGetChannel              Cmd = "GET_CHANNEL"
	CmdGetChannels             Cmd = "GET_CHANNELS"
	CmdSubscribe               Cmd = "SUBSCRIBE"
	CmdUnsubscribe             Cmd = "UNSUBSCRIBE"
	CmdSetUserVoiceSettings    Cmd = "SET_USER_VOICE_SETTINGS"
	CmdSelectVoiceChannel      Cmd = "SELECT_VOICE_CHANNEL"
	CmdGetSelectedVoiceChannel Cmd = "GET_SELECTED_VOICE_CHANNEL"
	CmdSelectTextChannel       Cmd = "SELECT_TEXT_CHANNEL"
	CmdGetVoiceSettings        Cmd = "GET_VOICE_SETTINGS"
	CmdSetVoiceSettings        Cmd = "SET_VOICE_SETTINGS"
	CmdSetCertifiedDevices     Cmd = "SET_CERTIFIED_DEVICES"
	CmdSetActivity             Cmd = "SET_ACTIVITY"
	CmdSendActivityJoinInvite  Cmd = "SEND_ACTIVITY_JOIN_INVITE"
	CmdCloseActivityRequest    Cmd = "CLOSE_ACTIVITY_REQUEST"
)

type CmdArgs interface {
	cmdArgs()
}

type Event struct {
	EventType
	CmdArgs
}

func UnmarshalEvent(eventType EventType, data MessageData) (Event, error) {
	var err error
	event := Event{
		EventType: eventType,
	}
	switch eventType {
	case EventReady:
		event.CmdArgs = nil

	case EventGuildStatus:
		if eventData, ok := data.(EventDataGuildStatus); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataGuildStatus")
		} else {
			event.CmdArgs = CmdArgsSubscribeGuild{GuildID: eventData.Guild.ID}
		}

	case EventGuildCreate:
		event.CmdArgs = nil

	case EventChannelCreate:
		event.CmdArgs = nil

	case EventVoiceChannelSelect:
		event.CmdArgs = nil

	case EventVoiceSettingsUpdate:
		event.CmdArgs = nil

	case EventVoiceStateCreate:
		if _, ok := data.(EventDataVoiceStateCreate); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataVoiceStateCreate")
		} else {
			// Discord doesn't return the channelID here so we can't link it back to specific events.
			event.CmdArgs = nil
		}

	case EventVoiceStateUpdate:
		if _, ok := data.(EventDataVoiceStateUpdate); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataVoiceStateUpdate")
		} else {
			// Discord doesn't return the channelID here so we can't link it back to specific events.
			event.CmdArgs = nil
		}

	case EventVoiceStateDelete:
		if _, ok := data.(EventDataVoiceStateDelete); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataVoiceStateDelete")
		} else {
			// Discord doesn't return the channelID here so we can't link it back to specific events.
			event.CmdArgs = nil
		}

	case EventVoiceConnectionStatus:
		event.CmdArgs = nil

	case EventMessageCreate:
		if eventData, ok := data.(EventDataMessageCreate); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataMessageCreate")
		} else {
			event.CmdArgs = CmdArgsSubscribeChannel{ChannelID: eventData.ChannelID}
		}

	case EventMessageUpdate:
		if eventData, ok := data.(EventDataMessageUpdate); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataMessageUpdate")
		} else {
			event.CmdArgs = CmdArgsSubscribeChannel{ChannelID: eventData.ChannelID}
		}

	case EventMessageDelete:
		if eventData, ok := data.(EventDataMessageDelete); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataMessageDelete")
		} else {
			event.CmdArgs = CmdArgsSubscribeChannel{ChannelID: eventData.ChannelID}
		}

	case EventSpeakingStart:
		if _, ok := data.(EventDataSpeakingStart); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataSpeakingStart")
		} else {
			// Discord doesn't return the channelID here so we can't link it back to specific events.
			event.CmdArgs = nil
		}

	case EventSpeakingStop:
		if _, ok := data.(EventDataSpeakingStop); !ok {
			err = errors.New("unable to cast CmdArgs to EventDataSpeakingStop")
		} else {
			// Discord doesn't return the channelID here so we can't link it back to specific events.
			event.CmdArgs = nil
		}

	case EventNotificationCreate:
		event.CmdArgs = nil

	case EventActivityJoin:
		event.CmdArgs = nil

	case EventActivitySpectate:
		event.CmdArgs = nil

	case EventActivityJoinRequest:
		event.CmdArgs = nil

	default:
		err = fmt.Errorf("unknown event: %s", eventType)
	}
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

type EventType string

const (
	EventReady                 EventType = "READY"
	EventError                 EventType = "ERROR"
	EventGuildStatus           EventType = "GUILD_STATUS"
	EventGuildCreate           EventType = "GUILD"
	EventChannelCreate         EventType = "CHANNEL_CREATE"
	EventVoiceChannelSelect    EventType = "VOICE_CHANNEL_SELECT"
	EventVoiceStateCreate      EventType = "VOICE_STATE_CREATE"
	EventVoiceStateUpdate      EventType = "VOICE_STATE_UPDATE"
	EventVoiceStateDelete      EventType = "VOICE_STATE_DELETE"
	EventVoiceSettingsUpdate   EventType = "VOICE_SETTINGS_UPDATE"
	EventVoiceConnectionStatus EventType = "VOICE_CONNECTION_STATUS"
	EventSpeakingStart         EventType = "SPEAKING_START"
	EventSpeakingStop          EventType = "SPEAKING_STOP"
	EventMessageCreate         EventType = "MESSAGE_CREATE"
	EventMessageUpdate         EventType = "MESSAGE_UPDATE"
	EventMessageDelete         EventType = "MESSAGE_DELETE"
	EventNotificationCreate    EventType = "NOTIFICATION_CREATE"
	EventActivityJoin          EventType = "ACTIVITY_JOIN"
	EventActivitySpectate      EventType = "ACTIVITY_SPECTATE"
	EventActivityJoinRequest   EventType = "ACTIVITY_JOIN_REQUEST"
)

type MessageData interface {
	messageData()
}
