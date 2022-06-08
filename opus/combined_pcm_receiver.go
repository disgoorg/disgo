package opus

import (
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

var (
	_ PCMReceiver            = (*CombinedPCMReceiver)(nil)
	_ voice.EventHandlerFunc = (*CombinedPCMReceiver)(nil)
)

func NewCombinedPCMReceiver(combinedPCMReceiver CombinedPCMReceiver2) *CombinedPCMReceiver {
	return &CombinedPCMReceiver{
		combinedPCMReceiver: combinedPCMReceiver,
		queue:               map[snowflake.ID][]*PCMPacket{},
	}
}

type CombinedPCMReceiver struct {
	combinedPCMReceiver CombinedPCMReceiver2
	queue               map[snowflake.ID][]*PCMPacket
}

func (r *CombinedPCMReceiver) HandlePCM(userID snowflake.ID, packet *PCMPacket) {

}

func (r *CombinedPCMReceiver) HandleEvent(opCode voice.GatewayOpcode, data voice.GatewayMessageData) {

}

type CombinedPCMPacket struct {
	SSRC       []uint32
	Sequence   []uint16
	Timestamps []uint32
	PCM        []int16
}

type CombinedPCMReceiver2 interface {
	HandlePCM(userIDs []snowflake.ID, packet *CombinedPCMPacket)
}
