package main

import (
	"fmt"
	"sync"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

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

func (h *echoHandler) HandleOpus(userID snowflake.ID, packet *voice.Packet) {
	h.queueMu.Lock()
	defer h.queueMu.Unlock()

	if len(h.queue) > 60 {
		println("dropping opus cause queue is full")
		return
	}

	h.queue = append(h.queue, packet.Opus)
}

func newEcho2(conn *voice.Connection) {
	conn.Speaking(voice.SpeakingFlagMicrophone)
	conn.UDP().Write(voice.SilenceFrames)
	for {
		packet, err := conn.UDP().ReadPacket()
		if err != nil {
			fmt.Printf("error while reading from reader: %s", err)
			continue
		}
		if _, err = conn.UDP().Write(packet.Opus); err != nil {
			fmt.Printf("error while writing to UDP: %s", err)
			continue
		}
	}
}
