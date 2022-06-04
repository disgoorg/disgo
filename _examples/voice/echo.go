package main

import (
	"fmt"

	"github.com/disgoorg/disgo/voice"
)

func newEcho2(conn *voice.Connection) {
	conn.Speaking(voice.SpeakingFlagMicrophone)
	conn.UDPConn().Write(voice.SilenceFrames)
	for {
		packet, err := conn.UDPConn().ReadPacket()
		if err != nil {
			fmt.Printf("error while reading from reader: %s", err)
			continue
		}
		if _, err = conn.UDPConn().Write(packet.Opus); err != nil {
			fmt.Printf("error while writing to UDP: %s", err)
			continue
		}
	}
}
