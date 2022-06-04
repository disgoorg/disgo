package main

import (
	"encoding/binary"
	"io"
	"sync"
	"time"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

var (
	_ voice.SendHandler    = (*echoHandler)(nil)
	_ voice.ReceiveHandler = (*echoHandler)(nil)
)

func newEchoHandler() *echoHandler {
	return &echoHandler{}
}

type echoHandler struct {
	queue   [][]byte
	queueMu sync.Mutex
}

func (h *echoHandler) ProvideOpus() ([]byte, error) {
	h.queueMu.Lock()
	defer h.queueMu.Unlock()

	if len(h.queue) == 0 {
		return nil, nil
	}

	var buff []byte
	buff, h.queue = h.queue[0], h.queue[1:]

	return buff, nil
}

func (h *echoHandler) HandleOpus(userID snowflake.ID, opus []byte) {
	if userID != 170939974227591168 {
		return
	}
	h.queueMu.Lock()
	defer h.queueMu.Unlock()

	if len(h.queue) > 10 {
		return
	}

	newBuff := make([]byte, len(opus))
	copy(newBuff, opus)
	h.queue = append(h.queue, newBuff)
}

func newReaderSendHandler(reader io.Reader) voice.SendHandler {
	return &audioSendHandler{
		reader: reader,
	}
}

type audioSendHandler struct {
	reader io.Reader
}

func (h *audioSendHandler) CanProvide() bool {
	return true
}

func (h *audioSendHandler) ProvideOpus() ([]byte, error) {
	var lenbuf [4]byte

	if _, err := h.reader.Read(lenbuf[:]); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err

	}

	buf := make([]byte, int64(binary.LittleEndian.Uint32(lenbuf[:])))
	if _, err := h.reader.Read(buf); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return buf, nil
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
