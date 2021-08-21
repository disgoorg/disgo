package discord

import (
	"strconv"
	"time"
)

var Epoch int64 = 1420070400000

// Snowflake is a general utility class around discord's IDs
type Snowflake string

// DeconstructedSnowflake contains the properties used by Discord for each ID
type DeconstructedSnowflake struct {
	Timestamp time.Time
	WorkerID  int64
	ProcessID int64
	Increment int64
}

// String returns the string representation of the Snowflake
func (s Snowflake) String() string {
	return string(s)
}

// Deconstruct returns DeconstructedSnowflake (https://discord.com/developers/docs/reference#snowflakes-snowflake-id-format-structure-left-to-right)
func (s Snowflake) Deconstruct() DeconstructedSnowflake {
	snowflake, _ := strconv.ParseInt(s.String(), 10, 64)
	return DeconstructedSnowflake{
		Timestamp: time.Unix(0, ((snowflake>>22)+Epoch)*1_000_000),
		WorkerID:  (snowflake & 0x3E0000) >> 17,
		ProcessID: (snowflake & 0x1F000) >> 12,
		Increment: snowflake & 0xFFF,
	}
}

// Timestamp returns a Time value of the snowflake
func (s Snowflake) Timestamp() time.Time {
	return s.Deconstruct().Timestamp
}

// NewSnowflake returns a new Snowflake based on the given timestamp
//goland:noinspection GoUnusedExportedFunction
func NewSnowflake(timestamp time.Time) Snowflake {
	return Snowflake(strconv.FormatInt(((timestamp.UnixNano()/1_000_000)-Epoch)<<22, 10))
}
