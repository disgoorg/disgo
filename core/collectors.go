package core

var _ Collectors = (*collectorsImpl)(nil)

type Collectors interface {
	NewMessageCollector(filter MessageFilter) (<-chan *Message, func())
	NewMessageReactionAddCollector(filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func())
	NewMessageReactionRemoveCollector(filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func())
	NewInteractionCollector(filter InteractionFilter) (<-chan Interaction, func())
	NewApplicationCommandInteractionCollector(filter ApplicationCommandInteractionFilter) (<-chan ApplicationCommandInteraction, func())
	NewSlashCommandCollector(filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func())
	NewMessageCommandCollector(filter MessageCommandInteractionFilter) (<-chan *MessageCommandInteraction, func())
	NewUserCommandCollector(filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func())
	NewComponentInteractionCollector(filter ComponentInteractionFilter) (<-chan ComponentInteraction, func())
	NewButtonClickCollector(filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func())
	NewSelectMenuSubmitCollector(filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func())
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
	NewApplicationCommandInteractionCollectorFunc func(bot *Bot, filter ApplicationCommandInteractionFilter) (<-chan ApplicationCommandInteraction, func())
	NewSlashCommandCollectorFunc                  func(bot *Bot, filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func())
	NewMessageCommandCollectorFunc                func(bot *Bot, filter MessageCommandInteractionFilter) (<-chan *MessageCommandInteraction, func())
	NewUserCommandCollectorFunc                   func(bot *Bot, filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func())
	NewComponentInteractionCollectorFunc          func(bot *Bot, filter ComponentInteractionFilter) (<-chan ComponentInteraction, func())
	NewButtonClickCollectorFunc                   func(bot *Bot, filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func())
	NewSelectMenuSubmitCollectorFunc              func(bot *Bot, filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func())
	NewAutocompleteCollectorFunc                  func(bot *Bot, filter AutocompleteInteractionFilter) (<-chan *AutocompleteInteraction, func())
}

type collectorsImpl struct {
	Bot *Bot
	CollectorsConfig
}

func (c *collectorsImpl) NewMessageCollector(filter MessageFilter) (<-chan *Message, func()) {
	return c.NewMessageCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewMessageCommandCollector(filter MessageCommandInteractionFilter) (<-chan *MessageCommandInteraction, func()) {
	return c.NewMessageCommandCollectorFunc(c.Bot, filter)
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

func (c *collectorsImpl) NewApplicationCommandInteractionCollector(filter ApplicationCommandInteractionFilter) (<-chan ApplicationCommandInteraction, func()) {
	return c.NewApplicationCommandInteractionCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewSlashCommandCollector(filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func()) {
	return c.NewSlashCommandCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewUserCommandCollector(filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func()) {
	return c.NewUserCommandCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewComponentInteractionCollector(filter ComponentInteractionFilter) (<-chan ComponentInteraction, func()) {
	return c.NewComponentInteractionCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewButtonClickCollector(filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func()) {
	return c.NewButtonClickCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewSelectMenuSubmitCollector(filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func()) {
	return c.NewSelectMenuSubmitCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewAutocompleteCollector(filter AutocompleteInteractionFilter) (<-chan *AutocompleteInteraction, func()) {
	return c.NewAutocompleteCollectorFunc(c.Bot, filter)
}
