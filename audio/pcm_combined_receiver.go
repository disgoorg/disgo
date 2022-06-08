package audio

import (
	"context"
	"encoding/binary"
	"io"
	"sync"
	"time"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

func NewPCMCombinerReceiver(pcmCombinedFrameReceiver PCMCombinedFrameReceiver, receiveUserFunc voice.ReceiveUserFunc) PCMFrameReceiver {
	if receiveUserFunc == nil {
		receiveUserFunc = func(userID snowflake.ID) bool {
			return true
		}
	}
	receiver := &pcmCombinerReceiver{
		pcmCombinedFrameReceiver: pcmCombinedFrameReceiver,
		receiveUserFunc:          receiveUserFunc,
	}
	go receiver.startCombinePackets()
	return receiver
}

type pcmCombinerReceiver struct {
	pcmCombinedFrameReceiver PCMCombinedFrameReceiver
	cancelFunc               context.CancelFunc
	receiveUserFunc          voice.ReceiveUserFunc
	queue                    map[snowflake.ID]*[]audioData
	queueMu                  sync.Mutex
}

func (r *pcmCombinerReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) {
	if r.receiveUserFunc == nil && !r.receiveUserFunc(userID) {
		return
	}
	r.queueMu.Lock()
	defer r.queueMu.Unlock()

	data := audioData{
		time:   time.Now().UnixMilli(),
		userID: userID,
		packet: packet,
	}
	if r.queue[userID] == nil {
		r.queue[userID] = &[]audioData{data}
	} else {
		*r.queue[userID] = append(*r.queue[userID], data)
	}
}

func (r *pcmCombinerReceiver) startCombinePackets() {
	lastFrameSent := time.Now().UnixMilli()
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		default:
			r.combinePackets()
			sleepTime := time.Duration(20 - (time.Now().UnixMilli() - lastFrameSent))
			if sleepTime > 0 {
				time.Sleep(sleepTime * time.Millisecond)
			}
			if time.Now().UnixMilli() < lastFrameSent+60 {
				lastFrameSent += 20
			} else {
				lastFrameSent = time.Now().UnixMilli()
			}
		}
	}
}

func (r *pcmCombinerReceiver) combinePackets() {
	r.queueMu.Lock()
	defer r.queueMu.Unlock()
	now := time.Now().UnixMilli()
	audioParts := make(map[snowflake.ID]audioData)
	var audioLen int
	for userID, packets := range r.queue {
		if len(*packets) == 0 {
			continue
		}

		var data *audioData
		*data, *packets = (*packets)[0], (*packets)[1:]
		for len(*packets) > 0 && now-data.time > 100 {
			*data, *packets = (*packets)[0], (*packets)[1:]
		}
		if data == nil {
			continue
		}
		audioParts[userID] = *data
		if len(data.packet.PCM) > audioLen {
			audioLen = len(data.packet.PCM)
		}
	}
	if len(audioParts) == 0 {
		return
	}
	combinedPacket := CombinedPCMPacket{
		SSRCs:      make([]uint32, len(audioParts)),
		Sequences:  make([]uint16, len(audioParts)),
		Timestamps: make([]uint32, len(audioParts)),
		PCM:        make([]int16, audioLen),
	}
	userIds := make([]snowflake.ID, 0, len(audioParts))
	for i, audio := range audioParts {
		combinedPacket.SSRCs[i] = audio.packet.SSRC
		combinedPacket.Sequences[i] = audio.packet.Sequence
		combinedPacket.Timestamps[i] = audio.packet.Timestamp
		userIds[i] = audio.userID

		for j := 0; j < len(audio.packet.PCM); j++ {
			newPCM := int32(combinedPacket.PCM[j]) + int32(audio.packet.PCM[j])
			if newPCM > 32767 {
				newPCM = 32767
			}
			if newPCM < -32768 {
				newPCM = -32768
			}
			combinedPacket.PCM[j] = int16(newPCM)
		}
	}
	r.pcmCombinedFrameReceiver.HandlePCM(userIds, &combinedPacket)
}

func (r *pcmCombinerReceiver) CleanupUser(userID snowflake.ID) {
	r.queueMu.Lock()
	defer r.queueMu.Unlock()
	delete(r.queue, userID)
}

func (r *pcmCombinerReceiver) Close() {
	r.cancelFunc()
	r.pcmCombinedFrameReceiver.Close()
}

type audioData struct {
	time   int64
	userID snowflake.ID
	packet *PCMPacket
}

type CombinedPCMPacket struct {
	SSRCs      []uint32
	Sequences  []uint16
	Timestamps []uint32
	PCM        []int16
}

type PCMCombinedFrameReceiver interface {
	HandlePCM(userIDs []snowflake.ID, packet *CombinedPCMPacket)
	Close()
}

func NewPCMCombinedStreamReceiver(w io.Writer) PCMCombinedFrameReceiver {
	return &pcmCombinedStreamReceiver{
		w: w,
	}
}

type pcmCombinedStreamReceiver struct {
	w io.Writer
}

func (p *pcmCombinedStreamReceiver) HandlePCM(_ []snowflake.ID, packet *CombinedPCMPacket) {
	_ = binary.Write(p.w, binary.LittleEndian, packet)
}

func (*pcmCombinedStreamReceiver) Close() {}
