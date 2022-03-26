package discord

import (
	"testing"
	"time"

	"github.com/disgoorg/disgo/json"
	"github.com/stretchr/testify/assert"
)

type timeJSON struct {
	Now *Time `json:"now"`
}

func TestTime_UnmarshalJSON(t *testing.T) {
	now := time.Date(1, 2, 3, 4, 5, 6, 0, time.UTC)
	format := now.Format(TimeFormat)

	var v timeJSON
	data := []byte("{\"now\":\"" + format + "\"}")
	err := json.Unmarshal(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, timeJSON{Now: &Time{Time: now}}, v)

	data = []byte("{\"now\":null}")
	err = json.Unmarshal(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, timeJSON{Now: nil}, v)
}

func TestTime_MarshalJSON(t *testing.T) {
	now := time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC)
	format := now.Format(TimeFormat)

	v := timeJSON{Now: &Time{Time: now}}
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "{\"now\":\""+format+"\"}", string(data))

	v = timeJSON{Now: nil}
	data, err = json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "{\"now\":null}", string(data))
}
