package gateway

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/disgoorg/json/v2"
	"github.com/gorilla/websocket"
	"github.com/klauspost/compress/zstd"
)

type CompressionType string

const (
	CompressionNone        CompressionType = "none"
	CompressionZlibPayload                 = "zlib-payload"
	CompressionZlibStream                  = "zlib-stream"
	CompressionZstdStream                  = "zstd-stream"
)

func (t CompressionType) isStreamCompression() bool {
	return t == CompressionZstdStream || t == CompressionZlibStream
}

func (t CompressionType) isPayloadCompression() bool {
	return t == CompressionZlibPayload
}

func (t CompressionType) newTransport(conn *websocket.Conn, logger *slog.Logger) transport {
	switch t {
	case CompressionZlibStream:
		return newZlibStreamTransport(conn, logger)
	case CompressionZstdStream:
		return newZstdStreamTransport(conn, logger)
	default:
		// zlibPayloadTransport supports both compressed (using zlib)
		// and uncompressed payloads
		//
		// The identify payload will state whether (some) payloads
		// will be compressed or not
		return newZlibPayloadTransport(conn, logger)
	}
}

var syncFlush = []byte{0x00, 0x00, 0xff, 0xff}

// [zstdStreamTransport]: for connections using zstd-stream compression
// [zlibStreamTransport]: for connections using zlib-stream compression
// [zlibPayloadTransport]: for connections using zlib payload compression or no compression
type transport interface {
	ReceiveMessage() (Message, error, error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

func parseMessage(r io.Reader, logger *slog.Logger) (Message, error) {
	if logger.Enabled(context.Background(), slog.LevelDebug) {
		buff := new(bytes.Buffer)
		r = io.TeeReader(r, buff)
		// This might seem a bit weird, but it's done such that it will print the
		// same data that the json decoder used, as not all data will end with an EOF
		defer func() {
			if buff.Len() > 0 {
				logger.Debug("received gateway message", slog.String("data", buff.String()))
			}
		}()
	}

	var message Message
	return message, json.NewDecoder(r).Decode(&message)
}

type zstdStreamTransport struct {
	conn      *websocket.Conn
	logger    *slog.Logger
	inflator  *zstd.Decoder
	pipeWrite *io.PipeWriter
}

func newZstdStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zstdStreamTransport {
	pRead, pWrite := io.Pipe()
	inflator, _ := zstd.NewReader(pRead, zstd.WithDecoderConcurrency(0))

	return &zstdStreamTransport{
		conn:      conn,
		logger:    logger,
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

func (t *zstdStreamTransport) ReceiveMessage() (Message, error, error) {
	mt, r, err := t.conn.NextReader()
	if err != nil {
		return Message{}, err, nil
	}

	if mt != websocket.BinaryMessage {
		return Message{}, fmt.Errorf("expected binary message, received %d", mt), nil
	}

	_, err = io.Copy(t.pipeWrite, r)
	if err != nil {
		return Message{}, err, nil
	}

	message, err := parseMessage(t.inflator, t.logger)
	return message, nil, err
}

type zlibStreamTransport struct {
	conn     *websocket.Conn
	logger   *slog.Logger
	inflator io.ReadCloser
	buffer   *bytes.Buffer
}

func newZlibStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zlibStreamTransport {
	buffer := new(bytes.Buffer)
	// The inflator cannot be created here because the creation will try
	// to consume the buffer, which is empty
	// see: https://github.com/golang/go/issues/58992

	return &zlibStreamTransport{
		conn:   conn,
		logger: logger,
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

func (t *zlibStreamTransport) ReceiveMessage() (Message, error, error) {
	t.buffer.Reset()

	for {
		mt, data, err := t.conn.ReadMessage()
		if err != nil {
			return Message{}, err, nil
		}
		if mt != websocket.BinaryMessage {
			return Message{}, fmt.Errorf("expected binary message, received %d", mt), nil
		}

		t.buffer.Write(data)

		if isFrameEnd(data) {
			break
		}
	}

	if t.inflator == nil {
		var err error

		t.inflator, err = zlib.NewReader(t.buffer)
		if err != nil {
			return Message{}, err, nil
		}
	}

	message, err := parseMessage(t.inflator, t.logger)
	return message, nil, err
}

type zlibPayloadTransport struct {
	conn   *websocket.Conn
	logger *slog.Logger
}

func newZlibPayloadTransport(conn *websocket.Conn, logger *slog.Logger) *zlibPayloadTransport {
	return &zlibPayloadTransport{
		conn:   conn,
		logger: logger,
	}
}

func (t *zlibPayloadTransport) Close() error {
	return t.conn.Close()
}

func (t *zlibPayloadTransport) WriteMessage(messageType int, data []byte) error {
	return t.conn.WriteMessage(messageType, data)
}

func (t *zlibPayloadTransport) ReceiveMessage() (Message, error, error) {
	mt, r, err := t.conn.NextReader()
	if err != nil {
		return Message{}, err, nil
	}

	if mt == websocket.BinaryMessage {
		reader, err := zlib.NewReader(r)
		if err != nil {
			return Message{}, fmt.Errorf("failed to decompress zlib: %w", err), nil
		}
		defer reader.Close()
		r = reader
	}

	message, err := parseMessage(r, t.logger)
	return message, nil, err
}
