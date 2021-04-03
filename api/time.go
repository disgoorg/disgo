package api

import (
	"bytes"
	"encoding/json"
	"time"
)

const timestampFormat = "2006-01-02T15:04:05.000000+00:00"

var emptyTime = []byte("\"\"")
var nullTime = []byte("null")

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	var ts string
	if !t.IsZero() {
		ts = t.String()
	}

	return []byte(`"` + ts + `"`), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var ts time.Time

	if bytes.Equal(emptyTime, data) || bytes.Equal(nullTime, data) {
		return nil
	}

	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}

	t.Time = ts
	return nil
}

func (t Time) String() string {
	return t.Format(timestampFormat)
}
