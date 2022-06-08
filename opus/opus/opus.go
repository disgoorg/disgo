package opus

/*
#cgo pkg-config: opus
#include <opus/opus.h>
*/
import "C"
import (
	"errors"
)

var (
	ErrEncoderNotInitialized     = errors.New("opus encoder not initialized")
	ErrEncoderAlreadyInitialized = errors.New("opus encoder already initialized")

	ErrDecoderNotInitialized     = errors.New("opus decoder not initialized")
	ErrDecoderAlreadyInitialized = errors.New("opus decoder already initialized")

	ErrRepacketizerNotInitialized     = errors.New("opus repacketizer not initialized")
	ErrRepacketizerAlreadyInitialized = errors.New("opus repacketizer already initialized")
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
