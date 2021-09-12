package discord

import (
	"bytes"
	"strconv"
	"time"
)

var (
	emptyJSONString = []byte(`""`)
	nullJSONString  = []byte(`null`)
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(emptyJSONString, data) || bytes.Equal(nullJSONString, data) {
		return nil
	}
	str, _ := strconv.Unquote(string(data))
	tt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	var parsed string
	if !t.IsZero() {
		parsed = t.Time.Format(time.RFC3339)
	}

	return []byte(`"` + parsed + `"`), nil
}
