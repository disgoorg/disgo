package main

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/disgoorg/disgo/voice"
)

func newReaderSendHandler(reader io.Reader) voice.OpusFrameProvider {
	return &audioSendHandler{
		reader: reader,
	}
}

type audioSendHandler struct {
	reader io.Reader
}

func (h *audioSendHandler) ProvideOpusFrame() []byte {
	var lenbuf [4]byte

	if _, err := h.reader.Read(lenbuf[:]); err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	buf := make([]byte, int64(binary.LittleEndian.Uint32(lenbuf[:])))
	if _, err := h.reader.Read(buf); err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	return buf
}

func writeOpus(w io.Writer, reader io.Reader) {
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenbuf [4]byte
	for {
		<-ticker.C
		_, err := io.ReadFull(reader, lenbuf[:])
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}

		// Read the integer
		framelen := int64(binary.LittleEndian.Uint32(lenbuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, reader, framelen)
		if err != nil && err != io.EOF {
			return
		}
	}
}
