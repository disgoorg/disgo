package voice

import (
	"testing"
	"time"

	"github.com/disgoorg/godave"
)

type notReadyDave struct {
	godave.Session
}

func (notReadyDave) Ready() bool {
	return false
}

func TestAudioReceiverWaitsWhileDaveIsNotReady(t *testing.T) {
	receiver := &defaultAudioReceiver{
		conn: &connImpl{dave: &notReadyDave{}},
	}

	start := time.Now()
	receiver.receive()

	if time.Since(start) < OpusFrameSizeMs*time.Millisecond/2 {
		t.Fatal("receive returned right away while DAVE is not ready, so the receiver loop spins")
	}
}
