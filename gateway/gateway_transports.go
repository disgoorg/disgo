package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/disgoorg/json/v2"
	"github.com/gorilla/websocket"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
)

// CompressionType defines the compression mechanism to use for a gateway connection
type CompressionType string

const (
	CompressionNone        CompressionType = ""
	CompressionZlibPayload CompressionType = "zlib-payload"
	CompressionZlibStream  CompressionType = "zlib-stream"
	CompressionZstdStream  CompressionType = "zstd-stream"
)

// IsStreamCompression returns whether the compression affects the whole stream
func (t CompressionType) IsStreamCompression() bool {
	return t == CompressionZstdStream || t == CompressionZlibStream
}

// IsPayloadCompression returns whether the compression affects only payloads
func (t CompressionType) IsPayloadCompression() bool {
	return t == CompressionZlibPayload
}

func (t CompressionType) String() string {
	if t == CompressionNone {
		return "none"
	}
	return string(t)
}

func newTransport(typ CompressionType, conn *websocket.Conn, logger *slog.Logger) transport {
	switch typ {
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

// transport is an abstraction over the underlying websocket connection to handle payload decompression.
//
//   - [zstdStreamTransport]: for connections using zstd-stream compression.
//   - [zlibStreamTransport]: for connections using zlib-stream compression.
//   - [zlibPayloadTransport]: for connections using zlib payload compression or no compression.
type transport interface {
	// ReceiveMessage returns the complete received [Message]. Message will be nil if it failed to parse it
	ReceiveMessage() (*Message, error)
	// WriteMessage writes a byte message to the underlying connection
	WriteMessage(message Message) error
	// WriteClose writes a close message to the underlying connection
	WriteClose(code int, message string) error
	// Close will free all resources and close the underlying connection
	Close() error
}

type baseTransport struct {
	conn   *websocket.Conn
	logger *slog.Logger
}

func (t *baseTransport) parseMessage(r io.Reader) (*Message, error) {
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
		return nil, err
	}

	return &message, nil
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

// pipeBuffer acts like a buffered pipe: it will return all the stored information, but
// avoid sending an EOF, as there can be more information still to come.
//
// It is important we don't embed bytes.Buffer, otherwise certain libs (ie, compress/zstd)
// will try to optimize the decompression using the .Bytes() and .Len() methods available
// in bytes.Buffer, making it impossible for us to prevent reporting an EOF.
type pipeBuffer struct{ buffer bytes.Buffer }

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read; err is always nil.
func (r *pipeBuffer) Read(p []byte) (int, error) {
	n, err := r.buffer.Read(p)
	if err == io.EOF {
		return n, nil
	}
	return n, err
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. If the buffer becomes too large, Write will panic with [ErrTooLarge].
func (r *pipeBuffer) Write(p []byte) {
	_, _ = r.buffer.Write(p)
}

// Reset resets the buffer to be empty, but it retains the underlying
// storage for use by future writes.
func (r *pipeBuffer) Reset() {
	r.buffer.Reset()
}

// zstdStreamTransport implements zstd-stream compression.
// See https://discord.com/developers/docs/events/gateway#zstdstream
type zstdStreamTransport struct {
	baseTransport

	inflator *zstd.Decoder
	buffer   *pipeBuffer
}

func newZstdStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zstdStreamTransport {
	return &zstdStreamTransport{
		baseTransport: baseTransport{
			conn:   conn,
			logger: logger,
		},
		buffer: new(pipeBuffer),
	}
}

func (t *zstdStreamTransport) ReceiveMessage() (*Message, error) {
	mt, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if mt != websocket.BinaryMessage {
		return nil, fmt.Errorf("expected binary message, received %d", mt)
	}

	t.buffer.Write(data)

	if t.inflator == nil {
		t.inflator, err = zstd.NewReader(t.buffer, zstd.WithDecoderConcurrency(1))
		if err != nil {
			return nil, err
		}
	}

	defer t.buffer.Reset()
	return t.parseMessage(t.inflator)
}

func (t *zstdStreamTransport) Close() error {
	t.buffer.Reset()
	if t.inflator != nil {
		t.inflator.Close()
	}
	return t.conn.Close()
}

// zlibStreamTransport implements zlib-stream compression.
// See https://discord.com/developers/docs/events/gateway#zlibstream
type zlibStreamTransport struct {
	baseTransport

	inflator io.ReadCloser
	buffer   *pipeBuffer
}

func newZlibStreamTransport(conn *websocket.Conn, logger *slog.Logger) *zlibStreamTransport {
	return &zlibStreamTransport{
		baseTransport: baseTransport{
			conn:   conn,
			logger: logger,
		},
		buffer: new(pipeBuffer),
	}
}

func isFrameEnd(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	return bytes.Equal(data[len(data)-4:], syncFlush)
}

func (t *zlibStreamTransport) ReceiveMessage() (*Message, error) {
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

	defer t.buffer.Reset()
	return t.parseMessage(t.inflator)
}

func (t *zlibStreamTransport) Close() error {
	t.buffer.Reset()
	if t.inflator != nil {
		_ = t.inflator.Close()
	}
	return t.conn.Close()
}

// zlibPayloadTransport implements both no compression and payload zlib compression.
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

	return t.parseMessage(r)
}

func (t *zlibPayloadTransport) Close() error {
	return t.conn.Close()
}
