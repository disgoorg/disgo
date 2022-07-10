package rpc

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"
)

const WSVersion = 1

func NewWSTransport(clientID snowflake.ID, origin string) (Transport, error) {
	var (
		conn *websocket.Conn
		err  error
	)
	for port := 6463; port < 6472; port++ {
		conn, err = openWS(clientID, origin, port)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return &wsTransport{
		conn: conn,
	}, nil
}

type wsTransport struct {
	conn *websocket.Conn
}

func (t *wsTransport) NextWriter() (io.WriteCloser, error) {
	return t.conn.NextWriter(websocket.TextMessage)
}

func (t *wsTransport) NextReader() (io.Reader, error) {
	mt, reader, err := t.conn.NextReader()
	if err != nil {
		return nil, err
	}
	if mt != websocket.TextMessage {
		return nil, errors.New("invalid message type")
	}
	return reader, nil
}

func (t *wsTransport) Close() error {
	return t.conn.Close()
}

func openWS(clientID snowflake.ID, origin string, port int) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d?v=%d&client_id=%d&encoding=json", port, WSVersion, clientID), http.Header{
		"Origin": []string{origin},
	})
	return conn, err
}
