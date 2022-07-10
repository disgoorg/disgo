package rpc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"

	"github.com/disgoorg/snowflake/v2"
)

type Handshake struct {
	V        int          `json:"v"`
	ClientID snowflake.ID `json:"client_id"`
}

type OpCode int32

const (
	OpCodeHandshake OpCode = iota
	OpCodeFrame
	OpCodeClose
	OpCodePing
	OpCodePong
)

func NewIPCTransport(clientID snowflake.ID, _ string) (Transport, error) {
	var (
		conn net.Conn
		err  error
	)
	for i := 0; i < 10; i++ {
		conn, err = openPipe(getDiscordIPCPath(i))
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	t := &ipcTransport{
		conn: conn,
		w:    bufio.NewWriter(conn),
	}

	if err = t.handshake(clientID); err != nil {
		return nil, err
	}

	return t, nil
}

type ipcTransport struct {
	conn net.Conn
	w    *bufio.Writer
}

func (t *ipcTransport) handshake(clientID snowflake.ID) error {
	w, err := t.nextWriter(OpCodeHandshake)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Close()
	}()

	return json.NewEncoder(w).Encode(Handshake{
		V:        Version,
		ClientID: clientID,
	})
}

func (t *ipcTransport) NextWriter() (io.WriteCloser, error) {
	return t.nextWriter(OpCodeFrame)
}

func (t *ipcTransport) nextWriter(opCode OpCode) (io.WriteCloser, error) {
	return &messageWriter{
		t:      t,
		opCode: opCode,
	}, nil
}

type messageWriter struct {
	t      *ipcTransport
	opCode OpCode
}

func (w *messageWriter) Write(p []byte) (int, error) {
	if err := binary.Write(w.t.w, binary.LittleEndian, w.opCode); err != nil {
		return 0, err
	}

	if err := binary.Write(w.t.w, binary.LittleEndian, int32(len(p))); err != nil {
		return 0, err
	}

	return w.t.w.Write(p)
}

func (w *messageWriter) Close() error {
	return w.t.w.Flush()
}

func (t *ipcTransport) NextReader() (io.Reader, error) {
	var opCode OpCode
	if err := binary.Read(t.conn, binary.LittleEndian, &opCode); err != nil {
		return nil, err
	}

	if opCode == OpCodeClose {
		_ = t.Close()
		return nil, net.ErrClosed
	}

	var length int32
	if err := binary.Read(t.conn, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := t.conn.Read(data); err != nil {
		return nil, err
	}

	if opCode == OpCodePing {
		if w, err := t.nextWriter(OpCodePong); err == nil {
			_, _ = w.Write(data)
			_ = w.Close()
		}
		return t.NextReader()
	}

	return bytes.NewReader(data), nil
}

func (t *ipcTransport) Close() error {
	if err := binary.Write(t.conn, binary.LittleEndian, OpCodeClose); err != nil {
		return err
	}
	return t.conn.Close()
}
