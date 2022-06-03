package main

import (
	"encoding/binary"
	"io"
	"sync"

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

	println("sending opus")
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

	println("received opus")
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

	if _, err := h.reader.Read(lenbuf[:]); err != nil {
		return nil, err
	}

	buf := make([]byte, int64(binary.LittleEndian.Uint32(lenbuf[:])))
	if _, err := h.reader.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
