package opus

/*
#cgo pkg-config: opus
#include <opus.h>
*/
import "C"
import "fmt"

var _ error = Error(0)

type Error int

func (e Error) Error() string {
	return fmt.Sprintf("opus: %s", string(C.GoString(C.opus_strerror(C.int(e)))))
}

const (
	ErrOK             Error = C.OPUS_OK
	ErrBadArg         Error = C.OPUS_BAD_ARG
	ErrBufferTooSmall Error = C.OPUS_BUFFER_TOO_SMALL
	ErrInternalError  Error = C.OPUS_INTERNAL_ERROR
	ErrInvalidPacket  Error = C.OPUS_INVALID_PACKET
	ErrUnimplemented  Error = C.OPUS_UNIMPLEMENTED
	ErrInvalidState   Error = C.OPUS_INVALID_STATE
	ErrAllocFail      Error = C.OPUS_ALLOC_FAIL
)
