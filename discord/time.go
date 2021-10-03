package discord

import (
	"bytes"
	"strconv"
	"time"

	"github.com/DisgoOrg/disgo/json"
)

const TimeFormat = "2006-01-02T15:04:05.000000+00:00"

var (
	emptyJSONString = []byte(`""`)
)

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)

type Time struct {
	time.Time
}


func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(emptyJSONString, data) {
		return nil
	}

	str, _ := strconv.Unquote(string(data))
	parsed, err := time.Parse(TimeFormat, str)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	var parsed string
	if !t.IsZero() {
		parsed = t.String()
	}

	return []byte(`"` + parsed + `"`), nil
}

func (t Time) String() string {
	return t.Format(TimeFormat)
}
