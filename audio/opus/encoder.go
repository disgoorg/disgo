package opus

/*
#cgo pkg-config: opus
#include <stdlib.h>
#include <opus/opus.h>
*/
import "C"

func NewEncoder(sampleRate int, channels int, application Application) (*Encoder, error) {
	var err C.int
	encoder := C.opus_encoder_create(C.opus_int32(sampleRate), C.int(channels), C.int(application), &err)
	if err != C.OPUS_OK {
		return nil, Error(err)
	}

	return &Encoder{
		encoder:  encoder,
		channels: channels,
	}, nil
}

type Encoder struct {
	encoder  *C.OpusEncoder
	channels int
}

func (e *Encoder) Init(sampleRate int, channels int, application Application) error {
	if e.encoder != nil {
		return ErrEncoderAlreadyInitialized
	}
	e.channels = channels
	if err := C.opus_encoder_init(e.encoder, C.opus_int32(sampleRate), C.int(channels), C.int(application)); err != C.OPUS_OK {
		return Error(err)
	}
	return nil
}

func (e *Encoder) Encode(pcm []int16, data []byte) (int, error) {
	if e.encoder == nil {
		return 0, ErrEncoderNotInitialized
	}
	n := C.opus_encode(e.encoder, (*C.opus_int16)(&pcm[0]), C.int(len(pcm)/e.channels), (*C.uchar)(&data[0]), C.opus_int32(cap(data)))
	if n < 0 {
		return 0, Error(n)
	}
	return int(n), nil
}

func (e *Encoder) EncodeFloat(pcm []float32, data []byte) (int, error) {
	if e.encoder == nil {
		return 0, ErrEncoderNotInitialized
	}
	n := C.opus_encode_float(e.encoder, (*C.float)(&pcm[0]), C.int(len(pcm)/e.channels), (*C.uchar)(&data[0]), C.opus_int32(cap(data)))
	if n < 0 {
		return 0, Error(n)
	}
	return int(n), nil
}

func (e *Encoder) Ctl(macro Macro[Encoder]) error {
	if e.encoder == nil {
		return ErrEncoderNotInitialized
	}
	if err := macro(e); err != C.OPUS_OK {
		return Error(err)
	}
	return nil
}

func (e *Encoder) Channels() int {
	return e.channels
}

func (e *Encoder) SampleRate() (int, error) {
	var sampleRate int
	if err := e.Ctl(GetEncoderSampleRate(&sampleRate)); err != nil {
		return 0, err
	}
	return sampleRate, nil
}

func (e *Encoder) Destroy() {
	if e.encoder == nil {
		return
	}
	C.opus_encoder_destroy(e.encoder)
}
