package gateway

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/klauspost/compress/zstd"
)

type CompressionType int

const (
	NoCompression CompressionType = iota
	ZlibPayloadCompression
	ZlibStreamCompression
	ZstdStreamCompression
)

var syncFlush = []byte{0x00, 0x00, 0xff, 0xff}

// [zstdStreamTransport]: for connections using zstd-stream compression
// [zlibStreamTransport]: for connections using zlib-stream compression
// [zlibPayloadTransport]: for connections using zlib payload compression or no compression
type transport interface {
	NextReader() (io.Reader, error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type zstdStreamTransport struct {
	conn      *websocket.Conn
	inflator  *zstd.Decoder
	pipeWrite *io.PipeWriter
}

func newZstdStreamTransport(conn *websocket.Conn) *zstdStreamTransport {
	pRead, pWrite := io.Pipe()
	inflator, _ := zstd.NewReader(pRead, zstd.WithDecoderConcurrency(0))

	return &zstdStreamTransport{
		conn:      conn,
		pipeWrite: pWrite,
		inflator:  inflator,
	}
}

func (t *zstdStreamTransport) Close() error {
	_ = t.pipeWrite.Close()
	t.inflator.Close()
	return t.conn.Close()
}

func (t *zstdStreamTransport) WriteMessage(messageType int, data []byte) error {
	return t.conn.WriteMessage(messageType, data)
}

func (t *zstdStreamTransport) NextReader() (io.Reader, error) {
	mt, r, err := t.conn.NextReader()
	if err != nil {
		return nil, err
	}

	if mt != websocket.BinaryMessage {
		return nil, fmt.Errorf("expected binary message, received %d", mt)
	}

	_, err = io.Copy(t.pipeWrite, r)
	if err != nil {
		return nil, err
	}

	return t.inflator, nil
}

type zlibStreamTransport struct {
	conn     *websocket.Conn
	inflator io.ReadCloser
	buffer   *bytes.Buffer
}

func newZlibStreamTransport(conn *websocket.Conn) *zlibStreamTransport {
	buffer := new(bytes.Buffer)
	// The inflator cannot be created here because the creation will try
	// to consume the buffer, which is empty
	// see: https://github.com/golang/go/issues/58992

	return &zlibStreamTransport{
		conn:   conn,
		buffer: buffer,
	}
}

func (t *zlibStreamTransport) Close() error {
	t.buffer.Reset()
	if t.inflator != nil {
		_ = t.inflator.Close()
	}
	return t.conn.Close()
}

func (t *zlibStreamTransport) WriteMessage(messageType int, data []byte) error {
	return t.conn.WriteMessage(messageType, data)
}

func isFrameEnd(data []byte) bool {
	return bytes.Equal(data[len(data)-4:], syncFlush)
}

func (t *zlibStreamTransport) NextReader() (io.Reader, error) {
	t.buffer.Reset()

	for {
		mt, data, err := t.conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		if mt != websocket.BinaryMessage {
			return nil, fmt.Errorf("expected binary message, received %d", mt)
		}

		t.buffer.Write(data)

		if isFrameEnd(data) {
			break
		}
	}

	var err error

	if t.inflator == nil {
		t.inflator, err = zlib.NewReader(t.buffer)
	}

	return t.inflator, err
}

type zlibPayloadTransport struct {
	conn *websocket.Conn
}

func newZlibPayloadTransport(conn *websocket.Conn) *zlibPayloadTransport {
	return &zlibPayloadTransport{
		conn: conn,
	}
}

func (t *zlibPayloadTransport) Close() error {
	return t.conn.Close()
}

func (t *zlibPayloadTransport) WriteMessage(messageType int, data []byte) error {
	return t.conn.WriteMessage(messageType, data)
}

func (t *zlibPayloadTransport) NextReader() (io.Reader, error) {
	mt, r, err := t.conn.NextReader()
	if err != nil {
		return nil, err
	}

	if mt == websocket.BinaryMessage {
		reader, err := zlib.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress zlib: %w", err)
		}
		defer reader.Close()
		r = reader
	}

	return r, nil
}
