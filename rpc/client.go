package rpc

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/insecurerandstr"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type EventHandleFunc func(event Event, data MessageData)

type Client interface {
	Logger() log.Logger
	ServerConfig() ServerConfig
	User() discord.User
	V() int

	Send(cmd Cmd, args CmdArgs, handler CommandHandler) error
	Close()
}

func NewClient(clientID snowflake.ID, eventHandleFunc EventHandleFunc) (Client, error) {
	conn, err := DialPipe()
	if err != nil {
		return nil, err
	}

	client := &clientImpl{
		logger:          log.Default(),
		conn:            conn,
		eventHandleFunc: eventHandleFunc,
		commandHandlers: map[string]internalHandler{},
		readyChan:       make(chan struct{}, 1),
	}

	buff := new(bytes.Buffer)
	if err = json.NewEncoder(buff).Encode(Handshake{
		V:        1,
		ClientID: clientID,
	}); err != nil {
		return nil, err
	}
	if err = client.send(OpCodeHandshake, buff); err != nil {
		return nil, err
	}

	go client.listen(conn)

	<-client.readyChan

	return client, nil
}

type clientImpl struct {
	logger          log.Logger
	conn            *Conn
	clientID        snowflake.ID
	eventHandleFunc EventHandleFunc
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

func (c *clientImpl) send(opCode OpCode, r io.Reader) error {
	writer, err := c.conn.NextWriter(opCode)
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
	c.logger.Debugf("Sending message: opCode: %d, data: %s", opCode, string(data))

	return err
}

func (c *clientImpl) Send(cmd Cmd, args CmdArgs, handler CommandHandler) error {
	nonce := insecurerandstr.RandStr(32)
	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(Message{
		Cmd:   cmd,
		Nonce: nonce,
		Args:  args,
	}); err != nil {
		return err
	}

	errChan := make(chan error, 1)

	c.commandHandlers[nonce] = internalHandler{
		handler: handler,
		errChan: errChan,
	}
	if err := c.send(OpCodeFrame, buff); err != nil {
		delete(c.commandHandlers, nonce)
		close(errChan)
		return err
	}
	return <-errChan
}

func (c *clientImpl) listen(conn *Conn) {
loop:
	for {
		println("reading...")
		opCode, reader, err := conn.NextReader()
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
		c.logger.Debugf("Received message: opCode: %d, data: %s", opCode, string(data))

		reader = bytes.NewReader(data)

		switch opCode {
		case OpCodePing:
			if err = c.send(OpCodePong, reader); err != nil {
				c.logger.Errorf("Error sending pong: %s", err)
				continue
			}

		case OpCodeFrame:
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
				c.eventHandleFunc(v.Event, v.Data)
				continue
			}
			if handler, ok := c.commandHandlers[v.Nonce]; ok {
				if v.Event == EventError {
					handler.errChan <- v.Data.(EventDataError)
				} else {
					handler.handler.Handle(v.Data)
					handler.errChan <- nil
				}
				close(handler.errChan)
				delete(c.commandHandlers, v.Nonce)
			} else {
				c.logger.Errorf("No handler for nonce: %s", v.Nonce)
			}

		case OpCodeClose:
			c.Close()
			break loop
		}
	}
}

func (c *clientImpl) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
