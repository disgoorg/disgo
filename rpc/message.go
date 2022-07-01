package rpc

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
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
	CmdSendActivityRequest     Cmd = "SEND_ACTIVITY_JOIN_INVITE"
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
