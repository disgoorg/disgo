package rpc

import (
	"fmt"

	"github.com/disgoorg/json"
)

type Message struct {
	Cmd  Cmd     `json:"cmd"`
	Args CmdArgs `json:"args,omitempty"`

	Event Event       `json:"evt,omitempty"`
	Data  MessageData `json:"data,omitempty"`

	Nonce string `json:"nonce,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v struct {
		Data json.RawMessage `json:"data"`
		message
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

type Event string

const (
	EventReady                 Event = "READY"
	EventError                 Event = "ERROR"
	EventGuildStatus           Event = "GUILD_STATUS"
	EventGuildCreate           Event = "GUILD"
	EventChannelCreate         Event = "CHANNEL_CREATE"
	EventVoiceChannelSelect    Event = "VOICE_CHANNEL_SELECT"
	EventVoiceStateCreate      Event = "VOICE_STATE_CREATE"
	EventVoiceStateUpdate      Event = "VOICE_STATE_UPDATE"
	EventVoiceStateDelete      Event = "VOICE_STATE_DELETE"
	EventVoiceSettingsUpdate   Event = "VOICE_SETTINGS_UPDATE"
	EventVoiceConnectionStatus Event = "VOICE_CONNECTION_STATUS"
	EventSpeakingStart         Event = "SPEAKING_START"
	EventSpeakingStop          Event = "SPEAKING_STOP"
	EventMessageCreate         Event = "MESSAGE_CREATE"
	EventMessageUpdate         Event = "MESSAGE_UPDATE"
	EventMessageDelete         Event = "MESSAGE_DELETE"
	EventNotificationCreate    Event = "NOTIFICATION_CREATE"
	EventActivityJoin          Event = "ACTIVITY_JOIN"
	EventActivitySpectate      Event = "ACTIVITY_SPECTATE"
	EventActivityJoinRequest   Event = "ACTIVITY_JOIN_REQUEST"
)

type MessageData interface {
	messageData()
}
