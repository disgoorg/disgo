package discord

import (
	"testing"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

func TestApplicationFlagsNew(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		flags   ApplicationFlags
		wantErr bool
	}{
		{
			name:    "flags only",
			payload: `{"id":"123","flags":4096}`,
			flags:   ApplicationFlagGatewayPresence,
		},
		{
			name:    "flags_new overrides flags",
			payload: `{"id":"123","flags":4096,"flags_new":"2147483648"}`,
			flags:   ApplicationFlags(1 << 31),
		},
		{
			name:    "flags_new alone",
			payload: `{"id":"123","flags_new":"2147483648"}`,
			flags:   ApplicationFlags(1 << 31),
		},
		{
			name:    "invalid flags_new",
			payload: `{"id":"123","flags":4096,"flags_new":"nope"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run("application/"+tt.name, func(t *testing.T) {
			var app Application
			err := json.Unmarshal([]byte(tt.payload), &app)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app.ID != snowflake.ID(123) {
				t.Errorf("id: got %d", app.ID)
			}
			if app.Flags != tt.flags {
				t.Errorf("flags: got %d, want %d", app.Flags, tt.flags)
			}
		})

		t.Run("partial/"+tt.name, func(t *testing.T) {
			var app PartialApplication
			err := json.Unmarshal([]byte(tt.payload), &app)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app != (PartialApplication{ID: snowflake.ID(123), Flags: tt.flags}) {
				t.Errorf("got %+v", app)
			}
		})
	}
}
