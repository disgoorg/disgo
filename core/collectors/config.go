package collectors

import "github.com/DisgoOrg/disgo/core"

var DefaultConfig = core.CollectorsConfig{
	NewMessageCollectorFunc:                       NewMessageCollector,
	NewMessageReactionAddCollectorFunc:            NewMessageReactionAddCollector,
	NewMessageReactionRemoveCollectorFunc:         NewMessageReactionRemoveCollector,
	NewInteractionCollectorFunc:                   NewInteractionCollector,
	NewApplicationCommandInteractionCollectorFunc: NewApplicationCommandInteractionCollector,
	NewComponentInteractionCollectorFunc:          NewComponentInteractionCollector,
	NewAutocompleteCollectorFunc:                  NewAutocompleteInteractionCollector,
}
