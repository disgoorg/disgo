package discord

import (
	"testing"

	"github.com/disgoorg/json/v2"
)

const (
	validButtonJSON      = `{"type": 2, "id": 2, "style": 1, "label": "click", "custom_id": "click"}`
	validTextDisplayJSON = `{"type": 10, "id": 3, "content": "hello"}`
	validThumbnailJSON   = `{"type": 11, "id": 4, "media": {"url": "https://example.com/image.png"}}`
)

func TestUnmarshalComponent_InvalidPosition(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		target  any
	}{
		{
			name:    "thumbnail inside container",
			payload: `{"type": 17, "id": 1, "components": [` + validThumbnailJSON + `]}`,
			target:  new(ContainerComponent),
		},
		{
			name:    "button as section child",
			payload: `{"type": 9, "id": 1, "components": [` + validButtonJSON + `], "accessory": ` + validThumbnailJSON + `}`,
			target:  new(SectionComponent),
		},
		{
			name:    "text display as section accessory",
			payload: `{"type": 9, "id": 1, "components": [` + validTextDisplayJSON + `], "accessory": ` + validTextDisplayJSON + `}`,
			target:  new(SectionComponent),
		},
		{
			name:    "text display inside action row",
			payload: `{"type": 1, "id": 1, "components": [` + validTextDisplayJSON + `]}`,
			target:  new(ActionRowComponent),
		},
		{
			name:    "button inside label",
			payload: `{"type": 18, "id": 1, "label": "label", "component": ` + validButtonJSON + `}`,
			target:  new(LabelComponent),
		},
		{
			name:    "button as message layout component",
			payload: `{"id": 1, "components": [` + validButtonJSON + `]}`,
			target:  new(Message),
		},
		{
			name:    "button as modal submit layout component",
			payload: `{"custom_id": "modal", "components": [` + validButtonJSON + `]}`,
			target:  new(ModalSubmitInteractionData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.payload), tt.target); err == nil {
				t.Errorf("expected error for invalid component position, got nil")
			}
		})
	}
}

func TestUnmarshalComponent_UnknownAndValidComponents(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		target  any
	}{
		{
			name:    "unknown component inside container",
			payload: `{"type": 17, "id": 1, "components": [{"type": 999, "id": 5}]}`,
			target:  new(ContainerComponent),
		},
		{
			name:    "valid section inside container",
			payload: `{"type": 17, "id": 1, "components": [{"type": 9, "id": 2, "components": [` + validTextDisplayJSON + `], "accessory": ` + validThumbnailJSON + `}]}`,
			target:  new(ContainerComponent),
		},
		{
			name:    "unknown component as message layout component",
			payload: `{"id": "1000000000000000001", "components": [{"type": 999, "id": 5}]}`,
			target:  new(Message),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.payload), tt.target); err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
