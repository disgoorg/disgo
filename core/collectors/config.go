package collectors

import "github.com/DisgoOrg/disgo/core"

var DefaultConfig = core.CollectorsConfig{
	NewMessageCollectorFunc:                       NewMessageCollector,
	NewMessageReactionAddCollectorFunc:            NewMessageReactionAddCollector,
	NewMessageReactionRemoveCollectorFunc:         NewMessageReactionRemoveCollector,
	NewInteractionCollectorFunc:                   NewInteractionCollector,
	NewApplicationCommandInteractionCollectorFunc: NewApplicationCommandInteractionCollector,
	NewSlashCommandCollectorFunc:                  NewSlashCommandCollector,
	NewMessageCommandCollectorFunc:                NewMessageCommandCollector,
	NewUserCommandCollectorFunc:                   NewUserCommandCollector,
	NewComponentInteractionCollectorFunc:          NewComponentInteractionCollector,
	NewButtonClickCollectorFunc:                   NewButtonClickCollector,
	NewSelectMenuSubmitCollectorFunc:              NewSelectMenuSubmitCollector,
	NewAutocompleteCollectorFunc:                  NewAutocompleteCollector,
}
