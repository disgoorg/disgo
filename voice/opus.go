package voice

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/disgoorg/snowflake/v2"
)

// NewOpusReader returns a new OpusFrameProvider that reads opus frames from the given io.Reader.
func NewOpusReader(r io.Reader) *OpusReader {
	return &OpusReader{
		r: r,
	}
}

// OpusReader is an OpusFrameProvider that reads opus frames from the given io.Reader.
// Each opus frame is prefixed with a 4 byte little endian uint32 that represents the length of the frame.
type OpusReader struct {
	r       io.Reader
	lenBuff [4]byte
	buff    [OpusFrameSizeBytes]byte
}

// ProvideOpusFrame reads the next opus frame from the underlying io.Reader.
func (h *OpusReader) ProvideOpusFrame() ([]byte, error) {
	_, err := h.r.Read(h.lenBuff[:])
	if err != nil {
		return nil, fmt.Errorf("error while reading opus frame length: %w", err)
	}

	frameLen := int64(binary.LittleEndian.Uint32(h.lenBuff[:]))
	actualLen, err := h.r.Read(h.buff[:frameLen])
	if err != nil {
		return nil, fmt.Errorf("error while reading opus frame: %w", err)
	}
	return h.buff[:actualLen], nil
}

// Close is a no-op.
func (*OpusReader) Close() {}

// NewOpusWriter returns a new OpusFrameReceiver that writes opus frames to the given io.Writer.
func NewOpusWriter(w io.Writer, userFilter UserFilterFunc) *OpusWriter {
	return &OpusWriter{
		w:          w,
		userFilter: userFilter,
	}
}

// OpusWriter is an OpusFrameReceiver that writes opus frames to the given io.Writer.
// Each opus frame is prefixed with a 4 byte little endian uint32 that represents the length of the frame.
type OpusWriter struct {
	w          io.Writer
	userFilter UserFilterFunc
}

// ReceiveOpusFrame writes the given opus frame to the underlying io.Writer.
func (r *OpusWriter) ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error {
	if r.userFilter != nil && !r.userFilter(userID) {
		return nil
	}
	if err := binary.Write(r.w, binary.LittleEndian, uint32(len(packet.Opus))); err != nil {
		return fmt.Errorf("error while writing opus frame length: %w", err)
	}
	if _, err := r.w.Write(packet.Opus); err != nil {
		return fmt.Errorf("error while writing opus frame: %w", err)
	}
	return nil
}

// CleanupUser is a no-op.
func (*OpusWriter) CleanupUser(_ snowflake.ID) {}

// Close is a no-op.
func (*OpusWriter) Close() {}
