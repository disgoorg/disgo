package opus

/*
#cgo pkg-config: opus
#include <opus/opus.h>
*/
import "C"
import (
	"errors"
)

const FrameSize = 20

var (
	ErrEncoderNotInitialized     = errors.New("audio encoder not initialized")
	ErrEncoderAlreadyInitialized = errors.New("audio encoder already initialized")

	ErrDecoderNotInitialized     = errors.New("audio decoder not initialized")
	ErrDecoderAlreadyInitialized = errors.New("audio decoder already initialized")

	ErrRepacketizerNotInitialized     = errors.New("audio repacketizer not initialized")
	ErrRepacketizerAlreadyInitialized = errors.New("audio repacketizer already initialized")
)

type Application int

const (
	ApplicationVoip               Application = C.OPUS_APPLICATION_VOIP
	ApplicationAudio              Application = C.OPUS_APPLICATION_AUDIO
	ApplicationRestrictedLowdelay Application = C.OPUS_APPLICATION_RESTRICTED_LOWDELAY
)

func Version() string {
	return C.GoString(C.opus_get_version_string())
}

func GetOutputBuffSize(rate int, channels int) int {
	return rate / 1000 * FrameSize * channels
}
