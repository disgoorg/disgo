package voice

import "time"

type SendHandler interface {
	CanProvide() bool
	ProvideOpus(frameLength time.Duration) ([]byte, error)
}

type ReceiveHandler interface {
	CanReceive() bool
	HandleOpus(opus []byte)
}
