package opus

/*
#cgo pkg-config: opus
#include <stdlib.h>
#include <opus/opus.h>
#include <opus/opus_defines.h>

int opus_set_complexity(OpusEncoder *st, opus_int32 complexity) {
	return opus_encoder_ctl(st, OPUS_SET_COMPLEXITY(complexity));
}
int opus_get_complexity(OpusEncoder *st, opus_int32 *complexity) {
	return opus_encoder_ctl(st, OPUS_GET_COMPLEXITY(complexity));
}

int opus_set_bitrate(OpusEncoder *st, opus_int32 bitrate) {
	return opus_encoder_ctl(st, OPUS_SET_BITRATE(bitrate));
}
int opus_get_bitrate(OpusEncoder *st, opus_int32 *bitrate) {
   return opus_encoder_ctl(st, OPUS_GET_BITRATE(bitrate));
}

int opus_set_vbr(OpusEncoder *st, opus_int32 vbr) {
	return opus_encoder_ctl(st, OPUS_SET_VBR(vbr));
}
int opus_get_vbr(OpusEncoder *st, opus_int32 *vbr) {
   return opus_encoder_ctl(st, OPUS_GET_VBR(vbr));
}

int opus_set_vbr_constraint(OpusEncoder *st, opus_int32 vbr_constraint) {
	return opus_encoder_ctl(st, OPUS_SET_VBR_CONSTRAINT(vbr_constraint));
}
int opus_get_vbr_constraint(OpusEncoder *st, opus_int32 *vbr_constraint) {
   return opus_encoder_ctl(st, OPUS_GET_VBR_CONSTRAINT(vbr_constraint));
}

int opus_set_force_channels(OpusEncoder *st, opus_int32 force_channels) {
	return opus_encoder_ctl(st, OPUS_SET_FORCE_CHANNELS(force_channels));
}
int opus_get_force_channels(OpusEncoder *st, opus_int32 *force_channels) {
   return opus_encoder_ctl(st, OPUS_GET_FORCE_CHANNELS(force_channels));
}

int opus_set_max_bandwidth(OpusEncoder *st, opus_int32 max_bandwidth) {
	return opus_encoder_ctl(st, OPUS_SET_MAX_BANDWIDTH(max_bandwidth));
}
int opus_get_max_bandwidth(OpusEncoder *st, opus_int32 *max_bandwidth) {
   return opus_encoder_ctl(st, OPUS_GET_MAX_BANDWIDTH(max_bandwidth));
}

int opus_set_bandwidth(OpusEncoder *st, opus_int32 bandwidth) {
	return opus_encoder_ctl(st, OPUS_SET_BANDWIDTH(bandwidth));
}
int opus_get_bandwidth(OpusEncoder *st, opus_int32 *bandwidth) {
   return opus_encoder_ctl(st, OPUS_GET_BANDWIDTH(bandwidth));
}

int opus_set_signal(OpusEncoder *st, opus_int32 signal) {
	return opus_encoder_ctl(st, OPUS_SET_SIGNAL(signal));
}
int opus_get_signal(OpusEncoder *st, opus_int32 *signal) {
   return opus_encoder_ctl(st, OPUS_GET_SIGNAL(signal));
}

int opus_set_application(OpusEncoder *st, opus_int32 application) {
	return opus_encoder_ctl(st, OPUS_SET_APPLICATION(application));
}
int opus_get_application(OpusEncoder *st, opus_int32 *application) {
   return opus_encoder_ctl(st, OPUS_GET_APPLICATION(application));
}

int opus_get_LOOKAHEAD(OpusEncoder *st, opus_int32 *lookahead) {
   return opus_encoder_ctl(st, OPUS_GET_LOOKAHEAD(lookahead));
}

int opus_set_inband_fec(OpusEncoder *st, opus_int32 fec) {
	return opus_encoder_ctl(st, OPUS_SET_INBAND_FEC(fec));
}
int opus_get_inband_fec(OpusEncoder *st, opus_int32 *fec) {
   return opus_encoder_ctl(st, OPUS_GET_INBAND_FEC(fec));
}

int opus_set_packet_loss_perc(OpusEncoder *st, opus_int32 perc) {
	return opus_encoder_ctl(st, OPUS_SET_PACKET_LOSS_PERC(perc));
}
int opus_get_packet_loss_perc(OpusEncoder *st, opus_int32 *perc) {
   return opus_encoder_ctl(st, OPUS_GET_PACKET_LOSS_PERC(perc));
}

int opus_set_dtx(OpusEncoder *st, opus_int32 dtx) {
	return opus_encoder_ctl(st, OPUS_SET_DTX(dtx));
}
int opus_get_dtx(OpusEncoder *st, opus_int32 *dtx) {
   return opus_encoder_ctl(st, OPUS_GET_DTX(dtx));
}

int opus_set_lsb_depth(OpusEncoder *st, opus_int32 depth) {
	return opus_encoder_ctl(st, OPUS_SET_LSB_DEPTH(depth));
}
int opus_get_lsb_depth(OpusEncoder *st, opus_int32 *depth) {
   return opus_encoder_ctl(st, OPUS_GET_LSB_DEPTH(depth));
}

int opus_set_expert_frame_duration(OpusEncoder *st, opus_int32 duration) {
	return opus_encoder_ctl(st, OPUS_SET_EXPERT_FRAME_DURATION(duration));
}
int opus_get_expert_frame_duration(OpusEncoder *st, opus_int32 *duration) {
   return opus_encoder_ctl(st, OPUS_GET_EXPERT_FRAME_DURATION(duration));
}

int opus_set_prediction_disabled(OpusEncoder *st, opus_int32 prediction) {
	return opus_encoder_ctl(st, OPUS_SET_PREDICTION_DISABLED(prediction));
}
int opus_get_prediction_disabled(OpusEncoder *st, opus_int32 *prediction) {
   return opus_encoder_ctl(st, OPUS_GET_PREDICTION_DISABLED(prediction));
}

int opus_encoder_get_sample_rate(OpusEncoder *st, opus_int32 *sample_rate) {
   return opus_encoder_ctl(st, OPUS_GET_SAMPLE_RATE(sample_rate));
}

int opus_decoder_get_sample_rate(OpusDecoder *st, opus_int32 *sample_rate) {
   return opus_decoder_ctl(st, OPUS_GET_SAMPLE_RATE(sample_rate));
}

int opus_set_phase_inversion_disabled(OpusEncoder *st, opus_int32 phase_inversion) {
	return opus_encoder_ctl(st, OPUS_SET_PHASE_INVERSION_DISABLED(phase_inversion));
}
int opus_get_phase_inversion_disabled(OpusEncoder *st, opus_int32 *phase_inversion) {
   return opus_encoder_ctl(st, OPUS_GET_PHASE_INVERSION_DISABLED(phase_inversion));
}

int opus_set_gain(OpusEncoder *st, opus_int32 gain) {
	return opus_encoder_ctl(st, OPUS_SET_GAIN(gain));
}
int opus_get_gain(OpusEncoder *st, opus_int32 *gain) {
   return opus_encoder_ctl(st, OPUS_GET_GAIN(gain));
}

int opus_get_last_packet_duration(OpusEncoder *st, opus_int32 *duration) {
   return opus_encoder_ctl(st, OPUS_GET_LAST_PACKET_DURATION(duration));
}

int opus_get_pitch(OpusEncoder *st, opus_int32 *pitch) {
   return opus_encoder_ctl(st, OPUS_GET_PITCH(pitch));
}

int opus_get_in_dtx(OpusEncoder *st, opus_int32 *in_dtx) {
   return opus_encoder_ctl(st, OPUS_GET_IN_DTX(in_dtx));
}

int opus_get_final_range(OpusEncoder *st, opus_uint32 *final_range) {
   return opus_encoder_ctl(st, OPUS_GET_FINAL_RANGE(final_range));
}

int opus_reset_state(OpusEncoder *st) {
	return opus_encoder_ctl(st, OPUS_RESET_STATE);
}
*/
import "C"
import "unsafe"

type Macro[T any] func(t *T) C.int

func SetBitrate(bitrate int) Macro[Encoder] {
	return func(e *Encoder) C.int {
		return C.opus_set_bitrate(e.encoder, C.int(bitrate))
	}
}

func GetEncoderSampleRate(sampleRate *int) Macro[Encoder] {
	return func(e *Encoder) C.int {
		return C.opus_encoder_get_sample_rate(e.encoder, (*C.int)(unsafe.Pointer(sampleRate)))
	}
}

func GetDecoderSamplerRte(sampleRate *int) Macro[Decoder] {
	return func(e *Decoder) C.int {
		return C.opus_decoder_get_sample_rate(e.decoder, (*C.int)(unsafe.Pointer(sampleRate)))
	}
}

func SetComplexity(complexity int) Macro[Encoder] {
	return func(e *Encoder) C.int {
		return C.opus_set_complexity(e.encoder, C.int(complexity))
	}
}

func SetForceChannels(channels int) Macro[Encoder] {
	return func(e *Encoder) C.int {
		return C.opus_set_force_channels(e.encoder, C.opus_int32(channels))
	}
}
