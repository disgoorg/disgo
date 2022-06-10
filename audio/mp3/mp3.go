package mp3

/*
#cgo pkg-config: libmpg123
#include <mpg123.h>
*/
import "C"

type Param C.mpg123_parms

const (
	ParamVerbose        Param = C.MPG123_VERBOSE
	ParamFlags          Param = C.MPG123_FLAGS
	ParamAddFlags       Param = C.MPG123_ADD_FLAGS
	ParamForceRate      Param = C.MPG123_FORCE_RATE
	ParamDownSample     Param = C.MPG123_DOWN_SAMPLE
	ParamRVA            Param = C.MPG123_RVA
	ParamDownspeed      Param = C.MPG123_DOWNSPEED
	ParamUpseed         Param = C.MPG123_UPSPEED
	ParamStartFrame     Param = C.MPG123_START_FRAME
	ParamDecodeFrames   Param = C.MPG123_DECODE_FRAMES
	ParamIcyInterval    Param = C.MPG123_ICY_INTERVAL
	ParamOutscale       Param = C.MPG123_OUTSCALE
	ParamTimeout        Param = C.MPG123_TIMEOUT
	ParamRemoveFlags    Param = C.MPG123_REMOVE_FLAGS
	ParamResyncLimit    Param = C.MPG123_RESYNC_LIMIT
	ParamIndexSize      Param = C.MPG123_INDEX_SIZE
	ParamPreframes      Param = C.MPG123_PREFRAMES
	ParamFeedpool       Param = C.MPG123_FEEDPOOL
	ParamFeedbuffer     Param = C.MPG123_FEEDBUFFER
	ParamFreeformatSize Param = C.MPG123_FREEFORMAT_SIZE
)

type ParamFlags C.mpg123_param_flags
