package mp3

/*
#cgo pkg-config: libmpg123
#include <stdlib.h>
#include <mpg123.h>
*/
import "C"
import (
	"io"
	"unsafe"
)

func CreateHandle() (*Handle, error) {
	var err C.int
	handle := C.mpg123_new(nil, &err)
	if err != C.MPG123_OK {
		return nil, Error(err)
	}

	err = C.mpg123_param(handle, C.MPG123_FORCE_RATE, 48000, 48000.0)
	if err != C.MPG123_OK {
		return nil, Error(err)
	}

	err = C.mpg123_param(handle, C.MPG123_ADD_FLAGS, C.MPG123_FORCE_STEREO, .0)
	if err != C.MPG123_OK {
		return nil, Error(err)
	}

	if err = C.mpg123_open_feed(handle); err != C.MPG123_OK {
		C.mpg123_delete(handle)
		return nil, Error(err)
	}

	return &Handle{
		handle: handle,
	}, nil
}

type Handle struct {
	handle *C.mpg123_handle
}

func (d *Handle) FormatNone() {
	C.mpg123_format_none(d.handle)
}

func (d *Handle) GetFormat() (int64, int, int) {
	var rate C.long
	var channels C.int
	var encoding C.int
	C.mpg123_getformat(d.handle, &rate, &channels, &encoding)
	return int64(rate), int(channels), int(encoding)
}

func (p *Handle) Param(param Param, ) {

}

func (d *Handle) Format(rate int64, channels int, encoding int) {
	C.mpg123_format(d.handle, C.long(rate), C.int(channels), C.int(encoding))
}

func (d *Handle) Feed(p []byte) (int, error) {
	if err := C.mpg123_feed(d.handle, (*C.uchar)(unsafe.Pointer(&p[0])), C.size_t(len(p))); err != C.MPG123_OK {
		return 0, Error(err)
	}
	return len(p), nil
}

func (d *Handle) Read(p []byte) (int, error) {
	var done C.size_t
	err := C.mpg123_read(d.handle, (unsafe.Pointer)(&p[0]), C.size_t(len(p)), &done)
	if err == C.MPG123_DONE {
		return int(done), io.EOF
	}
	if err != C.MPG123_OK {
		return int(done), Error(err)
	}
	return int(done), nil
}

func (d *Handle) Close() error {
	if err := C.mpg123_close(d.handle); err != C.MPG123_OK {
		return Error(err)
	}
	return nil
}
