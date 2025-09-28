package discord

import (
	"testing"

	"github.com/disgoorg/json/v2"
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
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}
	expected := `{"icon":"data:image/jpeg;base64,data"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	v = iconTest{Icon: nil}
	data, err = json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}
	expected = "{}"
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}
