package opus

/*
#cgo pkg-config: opus
#include <stdlib.h>
#include <opus/opus.h>
*/
import "C"

func NewDecoder(sampleRate int, channels int) (*Decoder, error) {
	var err C.int
	decoder := C.opus_decoder_create(C.opus_int32(sampleRate), C.int(channels), &err)
	if err != C.OPUS_OK {
		return nil, Error(err)
	}

	return &Decoder{
		decoder:  decoder,
		channels: channels,
	}, nil
}

type Decoder struct {
	decoder  *C.OpusDecoder
	channels int
}

func (e *Decoder) Init(sampleRate int, channels int) error {
	if e.decoder != nil {
		return ErrDecoderAlreadyInitialized
	}
	e.channels = channels
	if err := C.opus_decoder_init(e.decoder, C.opus_int32(sampleRate), C.int(channels)); err != C.OPUS_OK {
		return Error(err)
	}
	return nil
}

func (e *Decoder) Decode(data []byte, pcm []int16, decodeFec bool) (int, error) {
	if e.decoder == nil {
		return 0, ErrDecoderNotInitialized
	}
	var decodeFecCInt C.int
	if decodeFec {
		decodeFecCInt = C.int(1)
	}
	n := C.opus_decode(e.decoder, (*C.uchar)(&data[0]), C.opus_int32(len(data)), (*C.opus_int16)(&pcm[0]), C.int(cap(pcm)/e.channels), decodeFecCInt)
	if n < 0 {
		return 0, Error(n)
	}
	return int(n), nil
}

func (e *Decoder) DecodeFloat(data []byte, pcm []float32, decodeFec bool) (int, error) {
	if e.decoder == nil {
		return 0, ErrDecoderNotInitialized
	}
	var decodeFecCInt C.int
	if decodeFec {
		decodeFecCInt = C.int(1)
	}
	n := C.opus_decode_float(e.decoder, (*C.uchar)(&data[0]), C.opus_int32(len(data)), (*C.float)(&pcm[0]), C.int(cap(pcm)/e.channels), decodeFecCInt)
	if n < 0 {
		return 0, Error(n)
	}
	return int(n), nil
}

func (e *Decoder) Ctl(macro Macro[Decoder]) error {
	if e.decoder == nil {
		return ErrDecoderNotInitialized
	}
	if err := macro(e); err != C.OPUS_OK {
		return Error(err)
	}
	return nil
}

func (e *Decoder) Channels() int {
	return e.channels
}

func (e *Decoder) SampleRate() (int, error) {
	var sampleRate int
	if err := e.Ctl(GetDecoderSamplerRte(&sampleRate)); err != nil {
		return 0, err
	}
	return sampleRate, nil
}

func (e *Decoder) Destroy() {
	if e.decoder == nil {
		return
	}
	C.opus_decoder_destroy(e.decoder)
}
