package core

var _ Collectors = (*collectorsImpl)(nil)

type Collectors interface {
	NewMessageCollector(filter MessageFilter) (<-chan *Message, func())
	NewMessageReactionAddCollector(filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func())
	NewMessageReactionRemoveCollector(filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func())
	NewInteractionCollector(filter InteractionFilter) (<-chan Interaction, func())
	NewApplicationCommandInteractionCollector(filter ApplicationCommandInteractionFilter) (<-chan *ApplicationCommandInteraction, func())
	NewComponentInteractionCollector(filter ComponentInteractionFilter) (<-chan *ComponentInteraction, func())
	NewAutocompleteCollector(filter AutocompleteInteractionFilter) (<-chan *AutocompleteInteraction, func())
}

func NewCollectors(bot *Bot, config CollectorsConfig) Collectors {
	return &collectorsImpl{Bot: bot, CollectorsConfig: config}
}

type CollectorsConfig struct {
	NewMessageCollectorFunc                       func(bot *Bot, filter MessageFilter) (<-chan *Message, func())
	NewMessageReactionAddCollectorFunc            func(bot *Bot, filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func())
	NewMessageReactionRemoveCollectorFunc         func(bot *Bot, filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func())
	NewInteractionCollectorFunc                   func(bot *Bot, filter InteractionFilter) (<-chan Interaction, func())
	NewApplicationCommandInteractionCollectorFunc func(bot *Bot, filter ApplicationCommandInteractionFilter) (<-chan *ApplicationCommandInteraction, func())
	NewComponentInteractionCollectorFunc          func(bot *Bot, filter ComponentInteractionFilter) (<-chan *ComponentInteraction, func())
	NewAutocompleteCollectorFunc                  func(bot *Bot, filter AutocompleteInteractionFilter) (<-chan *AutocompleteInteraction, func())
}

type collectorsImpl struct {
	Bot *Bot
	CollectorsConfig
}

func (c *collectorsImpl) NewMessageCollector(filter MessageFilter) (<-chan *Message, func()) {
	return c.NewMessageCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewMessageReactionAddCollector(filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func()) {
	return c.NewMessageReactionAddCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewMessageReactionRemoveCollector(filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func()) {
	return c.NewMessageReactionRemoveCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewInteractionCollector(filter InteractionFilter) (<-chan Interaction, func()) {
	return c.NewInteractionCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewApplicationCommandInteractionCollector(filter ApplicationCommandInteractionFilter) (<-chan *ApplicationCommandInteraction, func()) {
	return c.NewApplicationCommandInteractionCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewComponentInteractionCollector(filter ComponentInteractionFilter) (<-chan *ComponentInteraction, func()) {
	return c.NewComponentInteractionCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewAutocompleteCollector(filter AutocompleteInteractionFilter) (<-chan *AutocompleteInteraction, func()) {
	return c.NewAutocompleteCollectorFunc(c.Bot, filter)
}
