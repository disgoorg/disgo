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

// transport is an abstraction over the underlying websocket connection to handle payload decompression
//
//   - [zstdStreamTransport]: for connections using zstd-stream compression
//   - [zlibStreamTransport]: for connections using zlib-stream compression
//   - [zlibPayloadTransport]: for connections using zlib payload compression or no compression
type transport interface {
	// ReceiveMessage returns the complete received [Message]. Message will be nil if it failed to parse it
	ReceiveMessage() (*Message, error)
	// WriteMessage writes a byte message to the underlying connection
	WriteMessage(message Message) error
	// WriteClose will free all resources and close the underlying connection
	WriteClose(code int, message string) error
	// Close will free all resources and close the underlying connection
	Close() error
}

type baseTransport struct {
	conn   *websocket.Conn
	logger *slog.Logger
}

func (t *baseTransport) parseMessage(r io.Reader) *Message {
	if t.logger.Enabled(context.Background(), slog.LevelDebug) {
		buff := new(bytes.Buffer)
		r = io.TeeReader(r, buff)
		// This might seem a bit weird, but it's done such that it will print the
		// same data that the json decoder used, as not all data will end with an EOF
		defer func() {
			if buff.Len() > 0 {
				t.logger.Debug("received gateway message", slog.String("data", buff.String()))
			}
		}()
	}

	var message Message
	err := json.NewDecoder(r).Decode(&message)
	if err != nil {
		t.logger.Error("error while parsing gateway message", slog.Any("err", err))
		return nil
	}

	return &message
}

func (t *baseTransport) WriteMessage(message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t.logger.Debug("sending gateway message", slog.String("data", string(data)))
	return t.conn.WriteMessage(websocket.TextMessage, data)
}

func (t *baseTransport) WriteClose(code int, message string) error {
	return t.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message))
}

// zstdStreamTransport implements zstd-stream compression
// See https://discord.com/developers/docs/events/gateway#zstdstream
type zstdStreamTransport struct {
	baseTransport

	inflator  *zstd.Decoder
	pipeWrite *io.PipeWriter
}

func newZstdStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zstdStreamTransport {
	pRead, pWrite := io.Pipe()
	inflator, _ := zstd.NewReader(pRead, zstd.WithDecoderConcurrency(0))

	return &zstdStreamTransport{
		baseTransport: baseTransport{
			conn:   conn,
			logger: logger,
		},
		pipeWrite: pWrite,
		inflator:  inflator,
	}
}

func (t *zstdStreamTransport) ReceiveMessage() (*Message, error) {
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

	return t.parseMessage(t.inflator), nil
}

func (t *zstdStreamTransport) Close() error {
	_ = t.pipeWrite.Close()
	t.inflator.Close()
	return t.conn.Close()
}

// zlibStreamTransport implements zlib-stream compression
// See https://discord.com/developers/docs/events/gateway#zlibstream
type zlibStreamTransport struct {
	baseTransport

	inflator io.ReadCloser
	buffer   *bytes.Buffer
}

func newZlibStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zlibStreamTransport {
	buffer := new(bytes.Buffer)
	// The inflator cannot be created here because the creation will try
	// to consume the buffer, which is empty
	// see: https://github.com/golang/go/issues/58992

	return &zlibStreamTransport{
		baseTransport: baseTransport{
			conn:   conn,
			logger: logger,
		},
		buffer: buffer,
	}
}

func isFrameEnd(data []byte) bool {
	return bytes.Equal(data[len(data)-4:], syncFlush)
}

func (t *zlibStreamTransport) ReceiveMessage() (*Message, error) {
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

	if t.inflator == nil {
		var err error

		t.inflator, err = zlib.NewReader(t.buffer)
		if err != nil {
			return nil, err
		}
	}

	return t.parseMessage(t.inflator), nil
}

func (t *zlibStreamTransport) Close() error {
	t.buffer.Reset()
	if t.inflator != nil {
		_ = t.inflator.Close()
	}
	return t.conn.Close()
}

// zlibPayloadTransport implements both no compression and payload zlib compression
// See https://discord.com/developers/docs/events/gateway#payload-compression
type zlibPayloadTransport struct {
	baseTransport
}

func newZlibPayloadTransport(conn *websocket.Conn, logger *slog.Logger) *zlibPayloadTransport {
	return &zlibPayloadTransport{
		baseTransport: baseTransport{
			conn:   conn,
			logger: logger,
		},
	}
}

func (t *zlibPayloadTransport) ReceiveMessage() (*Message, error) {
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

	return t.parseMessage(r), nil
}

func (t *zlibPayloadTransport) Close() error {
	return t.conn.Close()
}
