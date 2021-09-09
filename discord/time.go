package discord

import (
	"bytes"
	"time"
)

var emptyJSONString = []byte(`""`)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(emptyJSONString, data) {
		return nil
	}
	tt, err := time.Parse(time.RFC3339, string(data[1:len(data)-1]))
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
