package rpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/insecurerandstr"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var Version = 1

type TransportCreate func(clientID snowflake.ID, origin string) (Transport, error)

type Transport interface {
	NextWriter() (io.WriteCloser, error)
	NextReader() (io.Reader, error)
	Close() error
}

type Client interface {
	Logger() log.Logger
	ServerConfig() ServerConfig
	User() discord.User
	V() int
	Transport() Transport

	Subscribe(event Event, args CmdArgs, handler Handler) error
	Unsubscribe(event Event, args CmdArgs) error
	Send(message Message, handler Handler) error
	Close()
}

func NewClient(clientID snowflake.ID, opts ...ConfigOpt) (Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	client := &clientImpl{
		logger:          config.Logger,
		eventHandlers:   map[Event]Handler{},
		commandHandlers: map[string]internalHandler{},
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
	commandHandlers map[string]internalHandler

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

func (c *clientImpl) Subscribe(event Event, args CmdArgs, handler Handler) error {
	if _, ok := c.eventHandlers[event]; ok {
		return errors.New("event already subscribed")
	}
	c.eventHandlers[event] = handler
	return c.Send(Message{
		Cmd:   CmdSubscribe,
		Args:  args,
		Event: event,
	}, nil)
}

func (c *clientImpl) Unsubscribe(event Event, args CmdArgs) error {
	if _, ok := c.eventHandlers[event]; ok {
		delete(c.eventHandlers, event)
		return c.Send(Message{
			Cmd:   CmdUnsubscribe,
			Args:  args,
			Event: event,
		}, nil)
	}
	return nil
}

func (c *clientImpl) Send(message Message, handler Handler) error {
	nonce := insecurerandstr.RandStr(32)
	buff := new(bytes.Buffer)

	message.Nonce = nonce
	if err := json.NewEncoder(buff).Encode(message); err != nil {
		return err
	}

	errChan := make(chan error, 1)

	c.commandHandlers[nonce] = internalHandler{
		handler: handler,
		errChan: errChan,
	}
	if err := c.send(buff); err != nil {
		delete(c.commandHandlers, nonce)
		close(errChan)
		return err
	}
	return <-errChan
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
		if handler, ok := c.commandHandlers[v.Nonce]; ok {
			if v.Event == EventError {
				handler.errChan <- v.Data.(EventDataError)
			} else {
				if handler.handler != nil {
					handler.handler.Handle(v.Data)
				}
				handler.errChan <- nil
			}
			close(handler.errChan)
			delete(c.commandHandlers, v.Nonce)
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
