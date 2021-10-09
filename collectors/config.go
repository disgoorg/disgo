package collectors

import "github.com/DisgoOrg/disgo/core"

var DefaultConfig = core.CollectorsConfig{
	NewButtonClickCollectorFunc:           NewButtonClickCollector,
	NewMessageCollectorFunc:               NewMessageCollector,
	NewMessageCommandCollectorFunc:        NewMessageCommandCollector,
	NewMessageReactionAddCollectorFunc:    NewMessageReactionAddCollector,
	NewMessageReactionRemoveCollectorFunc: NewMessageReactionRemoveCollector,
	NewSelectMenuSubmitCollectorFunc:      NewSelectMenuSubmitCollector,
	NewSlashCommandCollectorFunc:          NewSlashCommandCollector,
	NewUserCommandCollectorFunc:           NewUserCommandCollector,
}
