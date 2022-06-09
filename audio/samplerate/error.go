package samplerate

/*
#cgo pkg-config: samplerate
#include <samplerate.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

type Error int

func (e Error) Error() string {
	return fmt.Sprintf("samplerate: %s", string(C.GoString(C.src_strerror(C.int(e)))))
}
