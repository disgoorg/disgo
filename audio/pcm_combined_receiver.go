package audio

import (
	"context"
	"encoding/binary"
	"io"
	"sync"
	"time"

	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

// NewPCMCombinerReceiver creates a new PCMFrameReceiver which combines multiple PCMPacket(s) into a single CombinedPCMPacket.
// You can process the CombinedPCMPacket by passing a PCMCombinedFrameReceiver.
func NewPCMCombinerReceiver(logger log.Logger, pcmCombinedFrameReceiver PCMCombinedFrameReceiver) PCMFrameReceiver {
	if logger == nil {
		logger = log.Default()
	}
	receiver := &pcmCombinerReceiver{
		logger:                   logger,
		pcmCombinedFrameReceiver: pcmCombinedFrameReceiver,
		queue:                    map[snowflake.ID]*[]audioData{},
	}
	go receiver.startCombinePackets()
	return receiver
}

type pcmCombinerReceiver struct {
	logger                   log.Logger
	pcmCombinedFrameReceiver PCMCombinedFrameReceiver
	cancelFunc               context.CancelFunc
	queue                    map[snowflake.ID]*[]audioData
	queueMu                  sync.Mutex
}

func (r *pcmCombinerReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) error {
	r.queueMu.Lock()
	defer r.queueMu.Unlock()

	pcm := make([]int16, len(packet.PCM))
	copy(pcm, packet.PCM)

	data := audioData{
		time:   time.Now().UnixMilli(),
		userID: userID,
		packet: &PCMPacket{
			SSRC:      packet.SSRC,
			Sequence:  packet.Sequence,
			Timestamp: packet.Timestamp,
			PCM:       pcm,
		},
	}

	if r.queue[userID] == nil {
		r.queue[userID] = &[]audioData{data}
	} else {
		*r.queue[userID] = append(*r.queue[userID], data)
	}
	return nil
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
			if err := r.combinePackets(); err != nil {
				r.logger.Error("Error combining pcm packets: ", err)
			}
			sleepTime := time.Duration(opus.FrameSize - (time.Now().UnixMilli() - lastFrameSent))
			if sleepTime > 0 {
				time.Sleep(sleepTime * time.Millisecond)
			}
			if time.Now().UnixMilli() < lastFrameSent+opus.FrameSize*2 {
				lastFrameSent += opus.FrameSize
			} else {
				lastFrameSent = time.Now().UnixMilli()
			}
		}
	}
}

func (r *pcmCombinerReceiver) combinePackets() error {
	r.queueMu.Lock()
	defer r.queueMu.Unlock()
	now := time.Now().UnixMilli()
	var audioParts []audioData
	var audioLen int
	for _, packets := range r.queue {
		if len(*packets) == 0 {
			continue
		}

		data := new(audioData)
		*data, *packets = (*packets)[0], (*packets)[1:]
		for len(*packets) > 0 && now-data.time > 100 {
			*data, *packets = (*packets)[0], (*packets)[1:]
		}
		if data == nil {
			continue
		}
		audioParts = append(audioParts, *data)
		if len(data.packet.PCM) > audioLen {
			audioLen = len(data.packet.PCM)
		}
	}
	if len(audioParts) == 0 {
		return nil
	}
	combinedPacket := &CombinedPCMPacket{
		Sequences:  make([]uint16, len(audioParts)),
		Timestamps: make([]uint32, len(audioParts)),
		SSRCs:      make([]uint32, len(audioParts)),
		PCM:        make([]int16, audioLen),
	}
	userIds := make([]snowflake.ID, len(audioParts))
	for i, audio := range audioParts {
		combinedPacket.Sequences[i] = audio.packet.Sequence
		combinedPacket.Timestamps[i] = audio.packet.Timestamp
		combinedPacket.SSRCs[i] = audio.packet.SSRC
		userIds[i] = audio.userID

		for j := 0; j < len(audio.packet.PCM); j++ {
			newPCM := int32(combinedPacket.PCM[j]) + int32(audio.packet.PCM[j]/int16(len(audioParts)))
			if newPCM > 32767 {
				newPCM = 32767
			}
			if newPCM < -32768 {
				newPCM = -32768
			}
			combinedPacket.PCM[j] = int16(newPCM)
		}
		i++
	}
	return r.pcmCombinedFrameReceiver.ReceiveCombinedPCMFrame(userIds, combinedPacket)
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

// CombinedPCMPacket is a PCMPacket which got created by combining multiple PCMPacket(s).
type CombinedPCMPacket struct {
	Sequences  []uint16
	Timestamps []uint32
	SSRCs      []uint32
	PCM        []int16
}

// PCMCombinedFrameReceiver is an interface for receiving PCMPacket(s) from multiple users as one CombinedPCMPacket.
type PCMCombinedFrameReceiver interface {
	// ReceiveCombinedPCMFrame is called when a new CombinedPCMPacket is received.
	ReceiveCombinedPCMFrame(userIDs []snowflake.ID, packet *CombinedPCMPacket) error

	// Close is called when the PCMCombinedFrameReceiver is no longer needed. It should close any open resources.
	Close()
}

// NewPCMCombinedStreamReceiver creates a new PCMCombinedFrameReceiver which writes the CombinedPCMPacket to the given io.Writer.
func NewPCMCombinedStreamReceiver(w io.Writer) PCMCombinedFrameReceiver {
	return &pcmCombinedStreamReceiver{
		w: w,
	}
}

type pcmCombinedStreamReceiver struct {
	w io.Writer
}

func (p *pcmCombinedStreamReceiver) ReceiveCombinedPCMFrame(_ []snowflake.ID, packet *CombinedPCMPacket) error {
	return binary.Write(p.w, binary.LittleEndian, packet.PCM)
}

func (*pcmCombinedStreamReceiver) Close() {}
