package audio

import (
	"sync"

	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

// NewPCMOpusReceiver creates a new voice.OpusFrameReceiver which receives Opus frames and decodes them into PCM frames. A new decoder is created for each user.
// You can pass your own *opus.Decoder by passing a decoderCreateFunc or nil to use the default Opus decoder(48000hz sample rate, 2 channels).
// You can filter users by passing a voice.ShouldReceiveUserFunc or nil to receive all users.
func NewPCMOpusReceiver(decoderCreateFunc func() *opus.Decoder, pcmFrameReceiver PCMFrameReceiver, receiveUserFunc voice.ShouldReceiveUserFunc) voice.OpusFrameReceiver {
	if decoderCreateFunc == nil {
		decoderCreateFunc = func() *opus.Decoder {
			decoder, err := opus.NewDecoder(48000, 2)
			if err != nil {
				panic("NewPCMOpusReceiver: " + err.Error())
			}
			return decoder
		}
	}
	if receiveUserFunc == nil {
		receiveUserFunc = func(userID snowflake.ID) bool {
			return true
		}
	}
	return &pcmOpusReceiver{
		receiveUserFunc:   receiveUserFunc,
		decoderCreateFunc: decoderCreateFunc,
		decoders:          map[snowflake.ID]*opus.Decoder{},
		pcmFrameReceiver:  pcmFrameReceiver,
	}
}

type pcmOpusReceiver struct {
	receiveUserFunc   voice.ShouldReceiveUserFunc
	decoderCreateFunc func() *opus.Decoder
	decoders          map[snowflake.ID]*opus.Decoder
	decodersMu        sync.Mutex
	pcmFrameReceiver  PCMFrameReceiver
	pcmBuff           [960 * 4]int16
}

func (r *pcmOpusReceiver) ReceiveOpusFrame(userID snowflake.ID, packet *voice.Packet) {
	if r.receiveUserFunc != nil && !r.receiveUserFunc(userID) {
		return
	}
	r.decodersMu.Lock()
	decoder, ok := r.decoders[userID]
	if !ok {
		decoder = r.decoderCreateFunc()
		r.decoders[userID] = decoder
	}
	r.decodersMu.Unlock()

	_, err := decoder.Decode(packet.Opus, r.pcmBuff[:], false)
	if err != nil {
		panic("ReceiveOpusFrame: " + err.Error())
		return
	}

	r.pcmFrameReceiver.ReceivePCMFrame(userID, &PCMPacket{
		SSRC:      packet.SSRC,
		Sequence:  packet.Sequence,
		Timestamp: packet.Timestamp,
		PCM:       r.pcmBuff[:],
	})
}

func (r *pcmOpusReceiver) CleanupUser(userID snowflake.ID) {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	decoder, ok := r.decoders[userID]
	if ok {
		decoder.Destroy()
		delete(r.decoders, userID)
	}
	r.pcmFrameReceiver.CleanupUser(userID)
}

func (r *pcmOpusReceiver) Close() {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	for _, decoder := range r.decoders {
		decoder.Destroy()
	}
	r.pcmFrameReceiver.Close()
}

// PCMPacket is a 20ms PCM frame with a ssrc, sequence and timestamp.
type PCMPacket struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	PCM       []int16
}
