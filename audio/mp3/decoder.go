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

var _ io.ReadWriteCloser = (*Decoder)(nil)

func CreateDecoder() (*Decoder, error) {
	var err C.int
	handle := C.mpg123_new(nil, &err)
	if err != C.MPG123_OK {
		return nil, Error(err)
	}

	return &Decoder{
		handle: handle,
	}, nil
}

type Decoder struct {
	handle *C.mpg123_handle
}

func (d *Decoder) FormatNone() {
	C.mpg123_format_none(d.handle)
}

func (d *Decoder) GetFormat() (int64, int, int) {
	var rate C.long
	var channels C.int
	var encoding C.int
	C.mpg123_getformat(d.handle, &rate, &channels, &encoding)
	return int64(rate), int(channels), int(encoding)
}

func (d *Decoder) Param(param Param, intValue int, floatValue float64) error {
	if err := C.mpg123_param(d.handle, C.int(param), C.long(intValue), C.double(floatValue)); err != C.MPG123_OK {
		return Error(err)
	}
	return nil
}

func (d *Decoder) Format(rate int64, channels int, encoding int) {
	C.mpg123_format(d.handle, C.long(rate), C.int(channels), C.int(encoding))
}

func (d *Decoder) Decode(in []byte, out []byte) (int, error) {
	var done C.size_t
	err := C.mpg123_decode(d.handle, (*C.uchar)(unsafe.Pointer(&in[0])), C.size_t(len(in)), (unsafe.Pointer)(&out[0]), C.size_t(cap(out)), &done)
	if err == C.MPG123_DONE {
		return int(done), io.EOF
	}
	if err != C.MPG123_OK {
		return int(done), Error(err)
	}
	return int(done), nil
}

func (d *Decoder) OpenFeed() error {
	if err := C.mpg123_open_feed(d.handle); err != C.MPG123_OK {
		return Error(err)
	}
	return nil
}

func (d *Decoder) Write(p []byte) (int, error) {
	if err := C.mpg123_feed(d.handle, (*C.uchar)(unsafe.Pointer(&p[0])), C.size_t(len(p))); err != C.MPG123_OK {
		return 0, Error(err)
	}
	return len(p), nil
}

func (d *Decoder) Read(p []byte) (int, error) {
	var done C.size_t
	err := C.mpg123_read(d.handle, (unsafe.Pointer)(&p[0]), C.size_t(cap(p)), &done)
	if err == C.MPG123_DONE {
		return int(done), io.EOF
	}
	if err == C.MPG123_NEW_FORMAT {
		return d.Read(p)
	}
	if err == C.MPG123_NEED_MORE {
		return 0, io.EOF
	}
	if err != C.MPG123_OK {
		return 0, Error(err)
	}
	return int(done), nil
}

func (d *Decoder) Close() error {
	if err := C.mpg123_close(d.handle); err != C.MPG123_OK {
		return Error(err)
	}
	return nil
}
