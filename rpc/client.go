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

func NewClient(clientID snowflake.ID, opts ...ConfigOpt) (Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	client := Client{
		Logger:          config.Logger,
		eventHandlers:   map[Event]Handler{},
		commandChannels: map[string]chan responseMessage{},
		readyChan:       make(chan struct{}, 1),
		clientId:        clientID,
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
	clientId        snowflake.ID

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

func (c *Client) send(r io.Reader) error {
	writer, err := c.Transport.NextWriter()
	if err != nil {
		return err
	}
	defer writer.Close()

	buff := new(bytes.Buffer)
	newWriter := io.MultiWriter(writer, buff)

	_, err = io.Copy(newWriter, r)
	if err != nil {
		return err
	}

	data, _ := io.ReadAll(buff)
	c.Logger.Tracef("Sending message: data: %s", string(data))

	return err
}

func (c *Client) Authorize(scopes []discord.OAuth2Scope, rpcToken, username string) (string, error) {
	args := CmdArgsAuthorize{
		ClientID: c.clientId,
		Scopes:   scopes,
	}

	if rpcToken != "" {
		args.RPCToken = rpcToken
	}

	if username != "" {
		args.Username = username
	}

	if res, err := c.Send(Message{
		Cmd:  CmdAuthorize,
		Args: args,
	}); err != nil {
		return "", err
	} else {
		return res.(CmdRsAuthorize).Code, nil
	}
}

func (c *Client) Authenticate(accessToken string) (CmdRsAuthenticate, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdAuthenticate,
		Args: CmdArgsAuthenticate{AccessToken: accessToken},
	}); err != nil {
		return CmdRsAuthenticate{}, err
	} else {
		return res.(CmdRsAuthenticate), nil
	}
}

func (c *Client) GetGuild(guildId snowflake.ID, timeout int) (CmdRsGetGuild, error) {
	args := CmdArgsGetGuild{
		GuildID: guildId,
	}

	if timeout != 0 {
		args.Timeout = timeout
	}

	if res, err := c.Send(Message{
		Cmd:  CmdGetGuild,
		Args: args,
	}); err != nil {
		return CmdRsGetGuild{}, err
	} else {
		return res.(CmdRsGetGuild), nil
	}
}

func (c *Client) GetGuilds() ([]PartialGuild, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetGuilds,
		Args: struct{ CmdArgs }{},
	}); err != nil {
		return nil, err
	} else {
		return res.(CmdRsGetGuilds).Guilds, nil
	}
}

func (c *Client) GetChannel(channelId snowflake.ID) (CmdRsGetChannel, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetChannel,
		Args: CmdArgsGetChannel{ChannelID: channelId},
	}); err != nil {
		return CmdRsGetChannel{}, err
	} else {
		return res.(CmdRsGetChannel), nil
	}
}

func (c *Client) GetChannels(guildId snowflake.ID) ([]PartialChannel, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetChannels,
		Args: CmdArgsGetChannels{GuildID: guildId},
	}); err != nil {
		return nil, err
	} else {
		return res.(CmdRsGetChannels).Channels, nil
	}
}

func (c *Client) GetVoiceSettings() (CmdRsGetVoiceSettings, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetVoiceSettings,
		Args: EmptyArgs{},
	}); err != nil {
		return CmdRsGetVoiceSettings{}, err
	} else {
		return res.(CmdRsGetVoiceSettings), nil
	}
}

func (c *Client) SetVoiceSettings(settings CmdArgsSetVoiceSettings) (CmdRsSetVoiceSettings, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdSetVoiceSettings,
		Args: settings,
	}); err != nil {
		return CmdRsSetVoiceSettings{}, err
	} else {
		return res.(CmdRsSetVoiceSettings), nil
	}
}

func (c *Client) SetUserVoiceSettings(settings CmdArgsSetUserVoiceSettings) (CmdRsSetUserVoiceSettings, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdSetUserVoiceSettings,
		Args: settings,
	}); err != nil {
		return CmdRsSetUserVoiceSettings{}, err
	} else {
		return res.(CmdRsSetUserVoiceSettings), nil
	}
}

func (c *Client) GetSelectedVoiceChannel() (*PartialChannel, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetSelectedVoiceChannel,
		Args: EmptyArgs{},
	}); err != nil {
		return &PartialChannel{}, err
	} else {
		return res.(CmdRsGetSelectedVoiceChannel).PartialChannel, nil
	}
}

func (c *Client) SelectVoiceChannel(channelID snowflake.ID, force bool, navigate bool) (*PartialChannel, error) {
	if res, err := c.Send(Message{
		Cmd: CmdSelectVoiceChannel,
		Args: CmdArgsSelectVoiceChannel{
			ChannelID: channelID,
			Force:     force,
			Navigate:  navigate,
		},
	}); err != nil {
		return &PartialChannel{}, err
	} else {
		return res.(CmdRsSelectVoiceChannel).PartialChannel, nil
	}
}

func (c *Client) SelectTextChannel(channelID *snowflake.ID) (*PartialChannel, error) {
	if res, err := c.Send(Message{
		Cmd: CmdSelectTextChannel,
		Args: CmdArgsSelectTextChannel{
			ChannelID: channelID,
		},
	}); err != nil {
		return &PartialChannel{}, err
	} else {
		return res.(CmdRsSelectTextChannel).PartialChannel, nil
	}
}

func (c *Client) SetActivity(PID int, activity discord.Activity) (CmdRsSetActivity, error) {
	if res, err := c.Send(Message{
		Cmd: CmdSetActivity,
		Args: CmdArgsSetActivity{
			PID:      PID,
			Activity: activity,
		},
	}); err != nil {
		return CmdRsSetActivity{}, err
	} else {
		log.Info(res)
		return res.(CmdRsSetActivity), err
	}
}

func (c *Client) Subscribe(event Event, args CmdArgs, handler Handler) error {
	if _, ok := c.eventHandlers[event]; ok {
		return errors.New("event already subscribed")
	}
	c.eventHandlers[event] = handler
	_, err := c.Send(Message{
		Cmd:   CmdSubscribe,
		Args:  args,
		Event: event,
	})
	return err
}

func (c *Client) Unsubscribe(event Event, args CmdArgs) error {
	if _, ok := c.eventHandlers[event]; ok {
		delete(c.eventHandlers, event)
		_, err := c.Send(Message{
			Cmd:   CmdUnsubscribe,
			Args:  args,
			Event: event,
		})
		return err
	}
	return nil
}

func (c *Client) Send(message Message) (MessageData, error) {
	nonce := insecurerandstr.RandStr(32)
	buff := new(bytes.Buffer)

	message.Nonce = nonce
	if err := json.NewEncoder(buff).Encode(message); err != nil {
		return nil, err
	}

	resChan := make(chan responseMessage, 1)

	c.commandChannels[nonce] = resChan

	if err := c.send(buff); err != nil {
		delete(c.commandChannels, nonce)
		close(resChan)
		return nil, err
	}

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
			if handler, ok := c.eventHandlers[v.Event]; ok {
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
