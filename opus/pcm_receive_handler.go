package opus

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/opus/opus"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

var _ voice.AudioReceiveHandler = (*PCMReceiveHandler)(nil)

func NewPCMReceiveHandler(decoderFunc func() *opus.Decoder, pcmReceiver PCMReceiver) *PCMReceiveHandler {
	return &PCMReceiveHandler{
		decoder:     decoder,
		pcmReceiver: pcmReceiver,
		pcmBuff:     make([]int16, 960),
		opusBuff:    make([]byte, 4000),
	}
}

type PCMReceiveHandler struct {
	decoderFunc func() *opus.Decoder
	decoders    map[snowflake.ID]*opus.Decoder
	pcmReceiver PCMReceiver
	pcmBuff     []int16
	opusBuff    []byte
}

func (h *PCMReceiveHandler) HandleOpus(userID snowflake.ID, packet *voice.Packet) {
	n, err := h.decoder.Decode(h.opusBuff, h.pcmBuff, true)
	if err != nil {
		return
	}

	r := bytes.NewReader(h.opusBuff[:n])
	if err = binary.Read(r, binary.LittleEndian, h.pcmBuff); err != nil {
		return
	}

	h.pcmReceiver.HandlePCM(userID, &PCMPacket{
		SSRC:      packet.SSRC,
		Sequence:  packet.Sequence,
		Timestamp: packet.Timestamp,
		PCM:       h.pcmBuff,
	})
}

type PCMPacket struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	PCM       []int16
}

type PCMReceiver interface {
	HandlePCM(userID snowflake.ID, packet *PCMPacket)
}

var _ PCMReceiver = (*PCMWriterReceiver)(nil)

func NewPCMWriterReceiver(w io.Writer) *PCMWriterReceiver {
	return &PCMWriterReceiver{
		w: w,
	}
}

type PCMWriterReceiver struct {
	w io.Writer
}

func (P PCMWriterReceiver) HandlePCM(_ snowflake.ID, packet *PCMPacket) {
	_ = binary.Write(P.w, binary.LittleEndian, packet.PCM)
}
