package voice

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/disgoorg/snowflake/v2"
)

// NewOpusReader returns a new OpusFrameProvider that reads opus frames from the given io.Reader.
func NewOpusReader(r io.Reader) OpusFrameProvider {
	return &opusReader{
		r:    r,
		buff: make([]byte, OpusStreamBuffSize),
	}
}

type opusReader struct {
	r       io.Reader
	lenBuff [4]byte
	buff    []byte
}

func (h *opusReader) ProvideOpusFrame() ([]byte, error) {
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

func (*opusReader) Close() {}

// NewOpusWriter returns a new OpusFrameReceiver that writes opus frames to the given io.Writer.
func NewOpusWriter(w io.Writer, userFilter UserFilterFunc) OpusFrameReceiver {
	return &opusWriter{
		w:          w,
		userFilter: userFilter,
	}
}

type opusWriter struct {
	w          io.Writer
	userFilter UserFilterFunc
}

func (r *opusWriter) ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error {
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

func (*opusWriter) CleanupUser(_ snowflake.ID) {}
func (*opusWriter) Close()                     {}
