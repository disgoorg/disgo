package samplerate

/*
#cgo pkg-config: samplerate
#include <samplerate.h>
#include <stdlib.h>
*/
import "C"

type ConverterType int

const (
	ConverterTypeSincBestQuality   ConverterType = C.SRC_SINC_BEST_QUALITY
	ConverterTypeSincMediumQuality ConverterType = C.SRC_SINC_MEDIUM_QUALITY
	ConverterTypeSincFastest       ConverterType = C.SRC_SINC_FASTEST
	ConverterTypeZeroOrderHold     ConverterType = C.SRC_ZERO_ORDER_HOLD
	ConverterTypeLinear            ConverterType = C.SRC_LINEAR
)
