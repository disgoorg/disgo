package samplerate

/*
#cgo pkg-config: samplerate
#include <samplerate.h>
#include <stdlib.h>
int bridge_src_process(SRC_STATE *state, float *data_in, float *data_out, long input_frames , long output_frames , int end_of_input, float src_ratio, long *input_frames_used, long *output_frames_gen) {
	SRC_DATA data;
	data.data_in = data_in;
	data.data_out = data_out;
	data.input_frames = input_frames;
	data.output_frames = output_frames;
	data.end_of_input = end_of_input;
	data.src_ratio = src_ratio;

	int err = src_process(state, &data);
	*input_frames_used = data.input_frames_used;
	*output_frames_gen = data.output_frames_gen;
	return err;
}
*/
import "C"

func CreateResampler(converterType ConverterType, channels int) *Resampler {
	var err C.int
	resampler := C.src_new(C.int(converterType), C.int(channels), &err)
	return &Resampler{
		resampler: resampler,
		channels:  channels,
	}
}

type Resampler struct {
	resampler *C.SRC_STATE
	channels  int
}

func (r *Resampler) Process(in []int16, out []int16, inputSampleRate int, outputSampleRate int, endOfInput int, inputFrames *int64, outputFrames *int64) error {
	inFloat := make([]float32, len(in))
	Int16ToFloat32Slice(in, inFloat)

	outFloat := make([]float32, cap(out))
	if err := r.ProcessFloat(inFloat, outFloat, inputSampleRate, outputSampleRate, endOfInput, inputFrames, outputFrames); err != nil {
		return err
	}
	Float32ToInt16Slice(outFloat, out)
	return nil
}

func (r *Resampler) ProcessFloat(in []float32, out []float32, inputSampleRate int, outputSampleRate int, endOfInput int, inputFrames *int64, outputFrames *int64) error {
	if err := C.bridge_src_process(r.resampler,
		(*C.float)(&in[0]),
		(*C.float)(&out[0]),
		C.long(len(in))/C.long(r.channels),
		C.long(cap(out))/C.long(r.channels),
		C.int(endOfInput),
		C.float(float64(outputSampleRate)/float64(inputSampleRate)),
		(*C.long)(inputFrames),
		(*C.long)(outputFrames),
	); err != 0 {
		return Error(err)
	}
	return nil
}

func (r *Resampler) Channels() int {
	return r.channels
}

func (r *Resampler) Destroy() {
	C.src_delete(r.resampler)
}

func Int16ToFloat32Slice(in []int16, out []float32) {
	C.src_short_to_float_array((*C.short)(&in[0]), (*C.float)(&out[0]), C.int(len(in)))
}

func Float32ToInt16Slice(in []float32, out []int16) {
	C.src_float_to_short_array((*C.float)(&in[0]), (*C.short)(&out[0]), C.int(len(in)))
}
