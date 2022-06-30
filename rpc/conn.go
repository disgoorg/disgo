package rpc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

func newConn(conn net.Conn) *Conn {
	return &Conn{
		Conn: conn,
		w:    bufio.NewWriter(conn),
	}
}

type Conn struct {
	net.Conn
	w *bufio.Writer
}

func (c *Conn) NextWriter(opCode OpCode) (io.WriteCloser, error) {
	return &messageWriter{
		conn:   c,
		opCode: opCode,
	}, nil
}

type messageWriter struct {
	conn   *Conn
	opCode OpCode
}

func (w *messageWriter) Write(p []byte) (int, error) {
	if err := binary.Write(w.conn.w, binary.LittleEndian, w.opCode); err != nil {
		return 0, err
	}

	if err := binary.Write(w.conn.w, binary.LittleEndian, int32(len(p))); err != nil {
		return 0, err
	}

	return w.conn.w.Write(p)
}

func (w *messageWriter) Close() error {
	return w.conn.w.Flush()
}

func (c *Conn) NextReader() (OpCode, io.Reader, error) {
	var opCode OpCode
	if err := binary.Read(c, binary.LittleEndian, &opCode); err != nil {
		return 0, nil, err
	}

	var length int32
	if err := binary.Read(c, binary.LittleEndian, &length); err != nil {
		return 0, nil, err
	}

	data := make([]byte, length)
	if _, err := c.Read(data); err != nil {
		return 0, nil, err
	}

	return opCode, bytes.NewReader(data), nil
}

func (c *Conn) Close() error {
	if err := binary.Write(c, binary.LittleEndian, OpCodeClose); err != nil {
		return err
	}
	return c.Conn.Close()
}
