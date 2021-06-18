package api

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// PatternTimestampFlag the regexp.Regexp to parse a Timestamp from a Message
var PatternTimestampFlag = regexp.MustCompile("<t:(?P<time>-?\\d{1,17})(?::(?P<format>[tTdDfFR]))?>")

// ErrNoTimestampMatch is returned when no valid Timestamp is found in the Message
var ErrNoTimestampMatch = errors.New("no matching timestamp found in string")

// TimestampFlag is used to determine how to display the Timestamp for the User in the client
type TimestampFlag string

const (
	// TimestampFlagShortTime formats time as 16:20
	TimestampFlagShortTime TimestampFlag = "t"

	// TimestampFlagLongTime formats time as 16:20:30
	TimestampFlagLongTime TimestampFlag = "T"

	// TimestampFlagShortDate formats time as 20/04/2021
	TimestampFlagShortDate TimestampFlag = "d"

	// TimestampFlagLongDate formats time as 20 April 2021
	TimestampFlagLongDate TimestampFlag = "D"

	// TimestampFlagShortDateTime formats time as 20 April 2021 16:20
	TimestampFlagShortDateTime TimestampFlag = "f"

	// TimestampFlagLongDateTime formats time as Tuesday, 20 April 2021 16:20
	TimestampFlagLongDateTime TimestampFlag = "F"

	// TimestampFlagRelative formats time as 2 months ago
	TimestampFlagRelative TimestampFlag = "R"
)

// FormatTime returns the time.Time formatted as markdown string
func (f TimestampFlag) FormatTime(time time.Time) string {
	return f.Format(time.Unix())
}

// Format returns the seconds formatted as markdown string
func (f TimestampFlag) Format(seconds int64) string {
	return fmt.Sprintf("<t:%d:%s>", seconds, f)
}

// ParseTimestamps parses all Timestamp(s) found in the provided string
func ParseTimestamps(str string, n int) ([]Timestamp, error) {
	matches := PatternTimestampFlag.FindAllStringSubmatch(str, n)
	if matches == nil {
		return nil, ErrNoTimestampMatch
	}

	timestamps := make([]Timestamp, len(matches))
	for i, match := range matches {
		unix, _ := strconv.Atoi(match[1])

		flag := TimestampFlagShortDateTime
		if len(match) > 2 {
			flag = TimestampFlag(match[2])
		}

		timestamps[i] = NewTimestampF(flag, time.Unix(int64(unix), 0))
	}

	return timestamps, nil
}

// ParseTimestamp parses the first Timestamp found in the provided string
func ParseTimestamp(str string) (*Timestamp, error) {
	timestamps, err := ParseTimestamps(str, 1)
	if err != nil {
		return nil, err
	}

	return &timestamps[0], nil
}

// NewTimestamp returns a new Timestamp with TimestampFlagShortDateTime & the given time.Time
func NewTimestamp(time time.Time) Timestamp {
	return NewTimestampF(TimestampFlagShortDateTime, time)
}

// NewTimestampF returns a new Timestamp with the given TimestampFlag & time.Time
func NewTimestampF(flag TimestampFlag, time time.Time) Timestamp {
	return Timestamp{
		TimestampFlag: flag,
		Time:          time,
	}
}

var _ fmt.Stringer = (*Timestamp)(nil)

// Timestamp represents a timestamp markdown object https://discord.com/developers/docs/reference#message-formatting
type Timestamp struct {
	time.Time
	TimestampFlag TimestampFlag
}

// String returns the Timestamp as markdown
func (t Timestamp) String() string {
	return t.Format()
}

// Format returns the Timestamp as markdown
func (t Timestamp) Format() string {
	return t.TimestampFlag.Format(t.Unix())
}

// FormatWith returns the Timestamp as markdown with the given TimestampFlag
func (t Timestamp) FormatWith(format TimestampFlag) string {
	return format.Format(t.Unix())
}
