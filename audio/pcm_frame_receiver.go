package audio

import (
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

// PCMFrameReceiver is an interface for receiving PCM frames.
type PCMFrameReceiver interface {
	// ReceivePCMFrame is called when a PCM frame is received.
	ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket)

	// CleanupUser is called when a user is disconnected. This should close any resources associated with the user.
	CleanupUser(userID snowflake.ID)

	// Close is called when the receiver is no longer needed. It should close any open resources.
	Close()
}

// NewPCMStreamReceiver creates a new PCMFrameReceiver which writes PCM frames to the given io.Writer.
// You can filter which users should be written by passing a voice.ShouldReceiveUserFunc.
func NewPCMStreamReceiver(w io.Writer, receiveUserFunc voice.ShouldReceiveUserFunc) PCMFrameReceiver {
	return &pcmStreamReceiver{
		w:               w,
		receiveUserFunc: receiveUserFunc,
	}
}

type pcmStreamReceiver struct {
	w               io.Writer
	receiveUserFunc voice.ShouldReceiveUserFunc
}

func (p *pcmStreamReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) {
	if p.receiveUserFunc == nil && !p.receiveUserFunc(userID) {
		return
	}
	_ = binary.Write(p.w, binary.LittleEndian, packet.PCM)
}

func (p *pcmStreamReceiver) CleanupUser(userID snowflake.ID) {}

func (*pcmStreamReceiver) Close() {}
