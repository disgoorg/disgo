package mp3

/*
#cgo pkg-config: libmpg123
#include <mpg123.h>
*/
import "C"
import "fmt"

type Error int

func (e Error) Error() string {
	return fmt.Sprintf("mp3: %s", string(C.GoString(C.mpg123_plain_strerror(C.int(e)))))
}

const (
	Done           Error = C.MPG123_DONE
	NewFormat      Error = C.MPG123_NEW_FORMAT
	NeedMore       Error = C.MPG123_NEED_MORE
	Err            Error = C.MPG123_ERR
	Ok             Error = C.MPG123_OK
	BadOutformat   Error = C.MPG123_BAD_OUTFORMAT
	BadChannel     Error = C.MPG123_BAD_CHANNEL
	BadRate        Error = C.MPG123_BAD_RATE
	Err16To08Table Error = C.MPG123_ERR_16TO8TABLE
)
