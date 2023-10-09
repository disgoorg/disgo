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

type Client interface {
	Logger() log.Logger
	ServerConfig() ServerConfig
	User() discord.User
	V() int
	Transport() Transport

	Authorize(args CmdArgsAuthorize) (string, error)
	Authenticate(args CmdArgsAuthenticate) (CmdRsAuthenticate, error)
	GetGuild(args CmdArgsGetGuild) (CmdRsGetGuild, error)
	GetGuilds() ([]PartialGuild, error)

	Subscribe(event Event, args CmdArgs, handler Handler) error
	Unsubscribe(event Event, args CmdArgs) error
	Send(message Message) (MessageData, error)
	Close()
}

func NewClient(clientID snowflake.ID, opts ...ConfigOpt) (Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	client := &clientImpl{
		logger:          config.Logger,
		eventHandlers:   map[Event]Handler{},
		commandChannels: map[string]chan responseMessage{},
		readyChan:       make(chan struct{}, 1),
	}

	if config.Transport == nil {
		var err error
		config.Transport, err = config.TransportCreate(clientID, config.Origin)
		if err != nil {
			return nil, err
		}
	}
	client.transport = config.Transport

	go client.listen(client.transport)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-client.readyChan:
	}

	return client, nil
}

type clientImpl struct {
	logger    log.Logger
	transport Transport

	eventHandlers   map[Event]Handler
	commandChannels map[string]chan responseMessage

	readyChan    chan struct{}
	user         discord.User
	serverConfig ServerConfig
	v            int
}

func (c *clientImpl) Logger() log.Logger {
	return c.logger
}

func (c *clientImpl) ServerConfig() ServerConfig {
	return c.serverConfig
}

func (c *clientImpl) User() discord.User {
	return c.user
}

func (c *clientImpl) V() int {
	return c.v
}

func (c *clientImpl) Transport() Transport {
	return c.transport
}

func (c *clientImpl) send(r io.Reader) error {
	writer, err := c.transport.NextWriter()
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
	c.logger.Tracef("Sending message: data: %s", string(data))

	return err
}

func (c *clientImpl) Authorize(args CmdArgsAuthorize) (string, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdAuthorize,
		Args: args,
	}); err != nil {
		return "", err
	} else {
		return res.(CmdRsAuthorize).Code, nil
	}
}

func (c *clientImpl) Authenticate(args CmdArgsAuthenticate) (CmdRsAuthenticate, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdAuthenticate,
		Args: args,
	}); err != nil {
		return CmdRsAuthenticate{}, err
	} else {
		return res.(CmdRsAuthenticate), nil
	}
}

func (c *clientImpl) GetGuild(args CmdArgsGetGuild) (CmdRsGetGuild, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetGuild,
		Args: args,
	}); err != nil {
		return CmdRsGetGuild{}, err
	} else {
		return res.(CmdRsGetGuild), nil
	}
}

func (c *clientImpl) GetGuilds() ([]PartialGuild, error) {
	if res, err := c.Send(Message{
		Cmd:  CmdGetGuilds,
		Args: struct{ CmdArgs }{},
	}); err != nil {
		return nil, err
	} else {
		return res.(CmdRsGetGuilds).Guilds, nil
	}
}

func (c *clientImpl) Subscribe(event Event, args CmdArgs, handler Handler) error {
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

func (c *clientImpl) Unsubscribe(event Event, args CmdArgs) error {
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

func (c *clientImpl) Send(message Message) (MessageData, error) {
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

func (c *clientImpl) listen(transport Transport) {
loop:
	for {
		reader, err := transport.NextReader()
		if errors.Is(err, net.ErrClosed) {
			c.logger.Error("Connection closed")
			break loop
		}
		if err != nil {
			c.logger.Errorf("Error reading message: %s", err)
			continue
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			c.logger.Errorf("Error reading message: %s", err)
			continue
		}
		c.logger.Tracef("Received message: data: %s", string(data))

		reader = bytes.NewReader(data)

		var v Message
		if err = json.NewDecoder(reader).Decode(&v); err != nil {
			c.logger.Errorf("failed to decode message: %s", err)
			continue
		}

		if v.Cmd == CmdDispatch {
			if d, ok := v.Data.(EventDataReady); ok {
				c.readyChan <- struct{}{}
				c.user = d.User
				c.serverConfig = d.Config
				c.v = d.V
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
			c.logger.Errorf("No handler for nonce: %s", v.Nonce)
		}
	}
}

func (c *clientImpl) Close() {
	if c.transport != nil {
		_ = c.transport.Close()
	}
}
