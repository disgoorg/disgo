package audio

import (
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

type PCMFrameReceiver interface {
	ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket)
	CleanupUser(userID snowflake.ID)
	Close()
}

func NewPCMStreamReceiver(w io.Writer, receiveUserFunc voice.ReceiveUserFunc) PCMFrameReceiver {
	return &pcmStreamReceiver{
		w:               w,
		receiveUserFunc: receiveUserFunc,
	}
}

type pcmStreamReceiver struct {
	w               io.Writer
	receiveUserFunc voice.ReceiveUserFunc
}

func (p *pcmStreamReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) {
	if p.receiveUserFunc == nil && !p.receiveUserFunc(userID) {
		return
	}
	_ = binary.Write(p.w, binary.LittleEndian, packet.PCM)
}

func (p *pcmStreamReceiver) CleanupUser(userID snowflake.ID) {}

func (*pcmStreamReceiver) Close() {}
