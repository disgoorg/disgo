package mp3

/*
#cgo pkg-config: libmpg123
#include <mpg123.h>
*/
import "C"

type Param C.int

const (
	Verbose        Param = C.MPG123_VERBOSE
	Flags          Param = C.MPG123_FLAGS
	AddFlags       Param = C.MPG123_ADD_FLAGS
	ForceRate      Param = C.MPG123_FORCE_RATE
	DownSample     Param = C.MPG123_DOWN_SAMPLE
	RVA            Param = C.MPG123_RVA
	Downspeed      Param = C.MPG123_DOWNSPEED
	Upseed         Param = C.MPG123_UPSPEED
	StartFrame     Param = C.MPG123_START_FRAME
	DecodeFrames   Param = C.MPG123_DECODE_FRAMES
	IcyInterval    Param = C.MPG123_ICY_INTERVAL
	Outscale       Param = C.MPG123_OUTSCALE
	Timeout        Param = C.MPG123_TIMEOUT
	RemoveFlags    Param = C.MPG123_REMOVE_FLAGS
	ResyncLimit    Param = C.MPG123_RESYNC_LIMIT
	IndexSize      Param = C.MPG123_INDEX_SIZE
	Preframes      Param = C.MPG123_PREFRAMES
	Feedpool       Param = C.MPG123_FEEDPOOL
	Feedbuffer     Param = C.MPG123_FEEDBUFFER
	FreeformatSize Param = C.MPG123_FREEFORMAT_SIZE
)

type ParamFlags C.int

const (
	ParamFlagForceMono          ParamFlags = C.MPG123_FORCE_MONO
	ParamFlagMonoLeft           ParamFlags = C.MPG123_MONO_LEFT
	ParamFlagMonoRight          ParamFlags = C.MPG123_MONO_RIGHT
	ParamFlagMonoMix            ParamFlags = C.MPG123_MONO_MIX
	ParamFlagForceStereo        ParamFlags = C.MPG123_FORCE_STEREO
	ParamFlagForce8Bit          ParamFlags = C.MPG123_FORCE_8BIT
	ParamFlagQuiet              ParamFlags = C.MPG123_QUIET
	ParamFlagGapless            ParamFlags = C.MPG123_GAPLESS
	ParamFlagNoResync           ParamFlags = C.MPG123_NO_RESYNC
	ParamFlagSeekbuffer         ParamFlags = C.MPG123_SEEKBUFFER
	ParamFlagFuzzy              ParamFlags = C.MPG123_FUZZY
	ParamFlagForceFloat         ParamFlags = C.MPG123_FORCE_FLOAT
	ParamFlagPlainID3text       ParamFlags = C.MPG123_PLAIN_ID3TEXT
	ParamFlagIgnoreStreamlength ParamFlags = C.MPG123_IGNORE_STREAMLENGTH
	ParamFlagSkipID3V           ParamFlags = C.MPG123_SKIP_ID3V2
	ParamFlagIgnoreInfoframe    ParamFlags = C.MPG123_IGNORE_INFOFRAME
	ParamFlagAutoResample       ParamFlags = C.MPG123_AUTO_RESAMPLE
	ParamFlagPicture            ParamFlags = C.MPG123_PICTURE
	ParamFlagNoPeelEnd          ParamFlags = C.MPG123_NO_PEEK_END
	ParamFlagForceSeekable      ParamFlags = C.MPG123_FORCE_SEEKABLE
	ParamFlagStoreRawID3        ParamFlags = C.MPG123_STORE_RAW_ID3
	ParamFlagForceEndian        ParamFlags = C.MPG123_FORCE_ENDIAN
	ParamFlagBigEndian          ParamFlags = C.MPG123_BIG_ENDIAN
	ParamFlagNoReadhead         ParamFlags = C.MPG123_NO_READAHEAD
	ParamFlagFloatFallback      ParamFlags = C.MPG123_FLOAT_FALLBACK
	ParamFlagNoFrankenstein     ParamFlags = C.MPG123_NO_FRANKENSTEIN
)
