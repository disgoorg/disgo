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
	conn, err := openWS(clientID, origin)
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

func openWS(clientID snowflake.ID, origin string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:6463?v=%d&client_id=%d&encoding=json", WSVersion, clientID), http.Header{
		"Origin": []string{origin},
	})
	return conn, err
}
