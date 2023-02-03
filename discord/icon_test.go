package discord

import (
	"testing"

	"github.com/disgoorg/json"
	"github.com/stretchr/testify/assert"
)

type iconTest struct {
	Icon *Icon `json:"icon,omitempty"`
}

func TestIcon_MarshalJSON(t *testing.T) {
	v := iconTest{
		Icon: &Icon{
			Type: IconTypeJPEG,
			Data: []byte("data"),
		},
	}

	data, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, `{"icon":"data:image/jpeg;base64,data"}`, string(data))

	v = iconTest{
		Icon: nil,
	}
	data, err = json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}
