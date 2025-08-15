package httpinteraction

import (
	"github.com/disgoorg/json/v2"

	"github.com/disgoorg/disgo/discord"
)

// EventInteractionCreate is the event payload when an interaction is created via Discord's Outgoing Webhooks
type EventInteractionCreate struct {
	discord.Interaction
}

func (e *EventInteractionCreate) UnmarshalJSON(data []byte) error {
	interaction, err := discord.UnmarshalInteraction(data)
	if err != nil {
		return err
	}
	e.Interaction = interaction
	return nil
}

func (e EventInteractionCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Interaction)
}
