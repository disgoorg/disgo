package rpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/disgoorg/disgo/internal/insecurerandstr"
	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var Version = 1

type TransportCreate func(clientID snowflake.ID, origin string) (Transport, error)

type Transport interface {
	NextWriter() (io.WriteCloser, error)
	NextReader() (io.Reader, error)
	Close() error
}

type responseMessage struct {
	data MessageData
	err  error
}

func New(clientID snowflake.ID, opts ...ConfigOpt) (Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	client := Client{
		Logger:          config.Logger,
		eventHandlers:   map[Event]Handler{},
		commandChannels: map[string]chan responseMessage{},
		readyChan:       make(chan struct{}, 1),
		clientID:        clientID,
	}

	if config.Transport == nil {
		var err error
		config.Transport, err = config.TransportCreate(clientID, config.Origin)
		if err != nil {
			return Client{}, err
		}
	}
	client.Transport = config.Transport

	return client, nil
}

type Client struct {
	Logger    log.Logger
	Transport Transport

	eventHandlers   map[Event]Handler
	commandChannels map[string]chan responseMessage
	clientID        snowflake.ID

	readyChan    chan struct{}
	User         discord.User
	ServerConfig ServerConfig
	V            int
}

func (c *Client) Open() error {
	go c.listen(c.Transport)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.readyChan:
	}
	return nil
}

func (c *Client) Authorize(scopes []discord.OAuth2Scope, rpcToken string, username string) (*CmdRsAuthorize, error) {
	args := CmdArgsAuthorize{
		ClientID: c.clientID,
		Scopes:   scopes,
		RPCToken: rpcToken,
		Username: username,
	}

	res, err := c.send(Message{
		Cmd:  CmdAuthorize,
		Args: args,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsAuthorize)), nil
}

func (c *Client) Authenticate(accessToken string) (*CmdRsAuthenticate, error) {
	res, err := c.send(Message{
		Cmd:  CmdAuthenticate,
		Args: CmdArgsAuthenticate{AccessToken: accessToken},
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsAuthenticate)), nil
}

func (c *Client) GetGuild(guildID snowflake.ID, timeout int) (*CmdRsGetGuild, error) {
	args := CmdArgsGetGuild{
		GuildID: guildID,
		Timeout: timeout,
	}

	res, err := c.send(Message{
		Cmd:  CmdGetGuild,
		Args: args,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsGetGuild)), nil
}

func (c *Client) GetGuilds() (*CmdRsGetGuilds, error) {
	res, err := c.send(Message{
		Cmd:  CmdGetGuilds,
		Args: nil,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsGetGuilds)), nil
}

func (c *Client) GetChannel(channelID snowflake.ID) (*CmdRsGetChannel, error) {
	res, err := c.send(Message{
		Cmd:  CmdGetChannel,
		Args: CmdArgsGetChannel{ChannelID: channelID},
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsGetChannel)), nil
}

func (c *Client) GetChannels(guildID snowflake.ID) (*CmdRsGetChannels, error) {
	res, err := c.send(Message{
		Cmd:  CmdGetChannels,
		Args: CmdArgsGetChannels{GuildID: guildID},
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsGetChannels)), nil
}

func (c *Client) GetVoiceSettings() (*CmdRsGetVoiceSettings, error) {
	res, err := c.send(Message{
		Cmd:  CmdGetVoiceSettings,
		Args: nil,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsGetVoiceSettings)), nil
}

func (c *Client) SetVoiceSettings(settings CmdArgsSetVoiceSettings) (*CmdRsSetVoiceSettings, error) {
	res, err := c.send(Message{
		Cmd:  CmdSetVoiceSettings,
		Args: settings,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsSetVoiceSettings)), nil
}

func (c *Client) SetUserVoiceSettings(settings CmdArgsSetUserVoiceSettings) (*CmdRsSetUserVoiceSettings, error) {
	res, err := c.send(Message{
		Cmd:  CmdSetUserVoiceSettings,
		Args: settings,
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsSetUserVoiceSettings)), nil
}

func (c *Client) GetSelectedVoiceChannel() (*CmdRsGetSelectedVoiceChannel, error) {
	res, err := c.send(Message{
		Cmd:  CmdGetSelectedVoiceChannel,
		Args: nil,
	})
	if err != nil {
		return nil, err
	}
	channel := res.(CmdRsGetSelectedVoiceChannel)
	if channel.ID == 0 {
		return nil, nil
	}
	return &channel, nil
}

func (c *Client) SelectVoiceChannel(channelID snowflake.ID, force bool, navigate bool) (*CmdRsSelectVoiceChannel, error) {
	res, err := c.send(Message{
		Cmd: CmdSelectVoiceChannel,
		Args: CmdArgsSelectVoiceChannel{
			ChannelID: channelID,
			Force:     force,
			Navigate:  navigate,
		},
	})
	if err != nil {
		return nil, err
	}
	channel := res.(CmdRsSelectVoiceChannel)
	if channel.ID == 0 {
		return nil, nil
	}
	return &channel, nil
}

func (c *Client) SelectTextChannel(channelID *snowflake.ID) (*CmdRsSelectTextChannel, error) {
	res, err := c.send(Message{
		Cmd: CmdSelectTextChannel,
		Args: CmdArgsSelectTextChannel{
			ChannelID: channelID,
		},
	})
	if err != nil {
		return nil, err
	}
	channel := res.(CmdRsSelectTextChannel)
	if channel.ID == 0 {
		return nil, nil
	}
	return &channel, nil
}

func (c *Client) SetActivity(pid int, activity discord.Activity) (*CmdRsSetActivity, error) {
	res, err := c.send(Message{
		Cmd: CmdSetActivity,
		Args: CmdArgsSetActivity{
			PID:      pid,
			Activity: activity,
		},
	})
	if err != nil {
		return nil, err
	}
	return json.Ptr(res.(CmdRsSetActivity)), nil
}

func (c *Client) SendActivityJoinInvite(userID snowflake.ID) error {
	_, err := c.send(Message{
		Cmd: CmdSendActivityJoinInvite,
		Args: CmdArgsSendActivityJoinInvite{
			UserID: userID,
		},
	})
	return err
}

func (c *Client) CloseActivityRequest(userID snowflake.ID) error {
	_, err := c.send(Message{
		Cmd: CmdCloseActivityRequest,
		Args: CmdArgsCloseActivityRequest{
			UserID: userID,
		},
	})
	return err
}

func (c *Client) SetCertifiedDevices(devices []CertifiedDevice) error {
	_, err := c.send(Message{
		Cmd: CmdSetCertifiedDevices,
		Args: CmdArgsSetCertifiedDevices{
			Devices: devices,
		},
	})
	return err
}

func (c *Client) SubscribeGuildStatus(guildID snowflake.ID, handler func(data EventDataGuildStatus)) error {
	return c.subscribe(EventGuildStatus,
		CmdArgsSubscribeGuild{GuildID: guildID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeGuildStatus(guildID snowflake.ID) error {
	return c.unsubscribe(EventGuildStatus, CmdArgsSubscribeGuild{GuildID: guildID})
}

func (c *Client) SubscribeGuildCreate(handler func(data EventDataGuildCreate)) error {
	return c.subscribe(EventGuildCreate,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeGuildCreate() error {
	return c.unsubscribe(EventGuildCreate, nil)
}

func (c *Client) SubscribeChannelCreate(handler func(data EventDataChannelCreate)) error {
	return c.subscribe(EventChannelCreate,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeChannelCreate() error {
	return c.unsubscribe(EventChannelCreate, nil)
}

func (c *Client) SubscribeVoiceChannelSelect(handler func(data EventDataVoiceChannelSelect)) error {
	return c.subscribe(EventVoiceChannelSelect,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceChannelSelect() error {
	return c.unsubscribe(EventVoiceChannelSelect, nil)
}

func (c *Client) SubscribeVoiceSettingsUpdate(handler func(data EventDataVoiceSettingsUpdate)) error {
	return c.subscribe(EventVoiceSettingsUpdate,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceSettingsUpdate() error {
	return c.unsubscribe(EventVoiceSettingsUpdate, nil)
}

func (c *Client) SubscribeVoiceStateCreate(channelID snowflake.ID, handler func(data EventDataVoiceStateCreate)) error {
	return c.subscribe(EventVoiceStateCreate,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceStateCreate(channelID snowflake.ID) error {
	return c.unsubscribe(EventVoiceStateCreate, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeVoiceStateUpdate(channelID snowflake.ID, handler func(data EventDataVoiceStateUpdate)) error {
	return c.subscribe(EventVoiceStateUpdate,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceStateUpdate(channelID snowflake.ID) error {
	return c.unsubscribe(EventVoiceStateUpdate, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeVoiceStateDelete(channelID snowflake.ID, handler func(data EventDataVoiceStateDelete)) error {
	return c.subscribe(EventVoiceStateDelete,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceStateDelete(channelID snowflake.ID) error {
	return c.unsubscribe(EventVoiceStateDelete, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeVoiceConnectionStatus(handler func(data EventDataVoiceConnectionStatus)) error {
	return c.subscribe(EventVoiceConnectionStatus,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeVoiceConnectionStatus() error {
	return c.unsubscribe(EventVoiceConnectionStatus, nil)
}

func (c *Client) SubscribeSpeakingStart(channelID snowflake.ID, handler func(data EventDataSpeakingStart)) error {
	return c.subscribe(EventSpeakingStart,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeSpeakingStart(channelID snowflake.ID) error {
	return c.unsubscribe(EventSpeakingStart, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeSpeakingStop(channelID snowflake.ID, handler func(data EventDataSpeakingStop)) error {
	return c.subscribe(EventSpeakingStop,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeSpeakingStop(channelID snowflake.ID) error {
	return c.unsubscribe(EventSpeakingStop, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeMessageCreate(channelID snowflake.ID, handler func(data EventDataMessageCreate)) error {
	return c.subscribe(EventMessageCreate,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeMessageCreate(channelID snowflake.ID) error {
	return c.unsubscribe(EventMessageCreate, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeMessageUpdate(channelID snowflake.ID, handler func(data EventDataMessageUpdate)) error {
	return c.subscribe(EventMessageUpdate,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeMessageUpdate(channelID snowflake.ID) error {
	return c.unsubscribe(EventMessageUpdate, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeMessageDelete(channelID snowflake.ID, handler func(data EventDataMessageDelete)) error {
	return c.subscribe(EventMessageDelete,
		CmdArgsSubscribeChannel{ChannelID: channelID},
		NewHandler(handler))
}

func (c *Client) UnsubscribeMessageDelete(channelID snowflake.ID) error {
	return c.unsubscribe(EventMessageDelete, CmdArgsSubscribeChannel{ChannelID: channelID})
}

func (c *Client) SubscribeNotificationCreate(handler func(data EventDataNotificationCreate)) error {
	return c.subscribe(EventNotificationCreate,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeNotificationCreate() error {
	return c.unsubscribe(EventNotificationCreate, nil)
}

func (c *Client) SubscribeActivityJoin(handler func(data EventDataActivityJoin)) error {
	return c.subscribe(EventActivityJoin,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeActivityJoin() error {
	return c.unsubscribe(EventActivityJoin, nil)
}

func (c *Client) SubscribeActivitySpectate(handler func(data EventDataActivitySpectate)) error {
	return c.subscribe(EventActivitySpectate,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeActivitySpectate() error {
	return c.unsubscribe(EventActivitySpectate, nil)
}

func (c *Client) SubscribeActivityJoinRequest(handler func(data EventDataActivityJoinRequest)) error {
	return c.subscribe(EventActivityJoinRequest,
		nil,
		NewHandler(handler))
}

func (c *Client) UnsubscribeActivityJoinRequest() error {
	return c.unsubscribe(EventActivityJoinRequest, nil)
}

func (c *Client) subscribe(event EventType, args CmdArgs, handler Handler) error {
	evt := Event{
		EventType: event,
		CmdArgs:   args,
	}

	// These events don't return channelID when fired. Set args empty for now to only allow single event.
	if event == EventVoiceStateCreate ||
		event == EventVoiceStateUpdate ||
		event == EventVoiceStateDelete ||
		event == EventSpeakingStart ||
		event == EventSpeakingStop {
		evt.CmdArgs = nil
	}

	if _, ok := c.eventHandlers[evt]; ok {
		return errors.New("event already subscribed")
	}
	c.eventHandlers[evt] = handler
	_, err := c.send(Message{
		Cmd:   CmdSubscribe,
		Args:  args,
		Event: event,
	})
	return err
}

func (c *Client) unsubscribe(event EventType, args CmdArgs) error {
	evt := Event{
		EventType: event,
		CmdArgs:   args,
	}

	// These events don't return channelID when fired. Set args empty for now to only allow single event.
	if event == EventVoiceStateCreate ||
		event == EventVoiceStateUpdate ||
		event == EventVoiceStateDelete ||
		event == EventSpeakingStart ||
		event == EventSpeakingStop {
		evt.CmdArgs = nil
	}

	if _, ok := c.eventHandlers[evt]; ok {
		delete(c.eventHandlers, evt)
		_, err := c.send(Message{
			Cmd:   CmdUnsubscribe,
			Args:  args,
			Event: event,
		})
		return err
	}
	return nil
}

func (c *Client) send(message Message) (MessageData, error) {
	nonce := insecurerandstr.RandStr(32)
	b := new(bytes.Buffer)

	message.Nonce = nonce
	if err := json.NewEncoder(b).Encode(message); err != nil {
		return nil, err
	}

	resChan := make(chan responseMessage, 1)

	c.commandChannels[nonce] = resChan

	writer, err := c.Transport.NextWriter()
	if err != nil {
		delete(c.commandChannels, nonce)
		close(resChan)
		return nil, err
	}

	buff := new(bytes.Buffer)
	newWriter := io.MultiWriter(writer, buff)

	_, err = io.Copy(newWriter, b)
	if err != nil {
		delete(c.commandChannels, nonce)
		close(resChan)
		writer.Close()
		return nil, err
	}

	data, _ := io.ReadAll(buff)
	c.Logger.Tracef("Sending message: data: %s", string(data))

	writer.Close()

	res := <-resChan

	return res.data, res.err
}

func (c *Client) listen(transport Transport) {
loop:
	for {
		reader, err := transport.NextReader()
		if errors.Is(err, net.ErrClosed) {
			c.Logger.Error("Connection closed")
			break loop
		}
		if err != nil {
			c.Logger.Errorf("Error reading message: %s", err)
			continue
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			c.Logger.Errorf("Error reading message: %s", err)
			continue
		}
		c.Logger.Tracef("Received message: data: %s", string(data))

		reader = bytes.NewReader(data)

		var v Message
		if err = json.NewDecoder(reader).Decode(&v); err != nil {
			c.Logger.Errorf("failed to decode message: %s", err)
			continue
		}

		if v.Cmd == CmdDispatch {
			if d, ok := v.Data.(EventDataReady); ok {
				c.readyChan <- struct{}{}
				c.User = d.User
				c.ServerConfig = d.Config
				c.V = d.V
			}
			event, err := UnmarshalEvent(v.Event, v.Data)
			if err != nil {
				c.Logger.Errorf("failed to build Event for eventType %s. %s", v.Event, err)
				continue
			}
			if handler, ok := c.eventHandlers[event]; ok {
				handler.Handle(v.Data)
			}
			continue
		}

		if handler, ok := c.commandChannels[v.Nonce]; ok {
			res := responseMessage{
				data: nil,
				err:  nil,
			}
			if v.Event == EventError {
				res.err = v.Data.(EventDataError)
			} else {
				res.data = v.Data
			}

			handler <- res
			close(handler)
			delete(c.commandChannels, v.Nonce)
		} else {
			c.Logger.Errorf("No handler for nonce: %s", v.Nonce)
		}
	}
}

func (c *Client) Close() {
	if c.Transport != nil {
		_ = c.Transport.Close()
	}
}
