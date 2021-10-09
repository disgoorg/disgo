package core

var _ Collectors = (*collectorsImpl)(nil)

type Collectors interface {
	NewButtonClickCollector(filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func())
	NewMessageCollector(filter MessageFilter) (<-chan *Message, func())
	NewMessageCommandCollector(filter MessageCommandInteractionFilter) (<-chan *MessageCommandInteraction, func())
	NewMessageReactionAddCollector(filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func())
	NewMessageReactionRemoveCollector(filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func())
	NewSelectMenuSubmitCollector(filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func())
	NewSlashCommandCollector(filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func())
	NewUserCommandCollector(filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func())
}

func NewCollectors(bot *Bot, config CollectorsConfig) Collectors {
	return &collectorsImpl{Bot: bot, CollectorsConfig: config}
}

type CollectorsConfig struct {
	NewButtonClickCollectorFunc           func(bot *Bot, filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func())
	NewMessageCollectorFunc               func(bot *Bot, filter MessageFilter) (<-chan *Message, func())
	NewMessageCommandCollectorFunc        func(bot *Bot, filter MessageCommandInteractionFilter) (<-chan *MessageCommandInteraction, func())
	NewMessageReactionAddCollectorFunc    func(bot *Bot, filter MessageReactionAddFilter) (<-chan *MessageReactionAdd, func())
	NewMessageReactionRemoveCollectorFunc func(bot *Bot, filter MessageReactionRemoveFilter) (<-chan *MessageReactionRemove, func())
	NewSelectMenuSubmitCollectorFunc      func(bot *Bot, filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func())
	NewSlashCommandCollectorFunc          func(bot *Bot, filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func())
	NewUserCommandCollectorFunc           func(bot *Bot, filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func())
}

type collectorsImpl struct {
	Bot *Bot
	CollectorsConfig
}

func (c *collectorsImpl) NewButtonClickCollector(filter ButtonInteractionFilter) (<-chan *ButtonInteraction, func()) {
	return c.NewButtonClickCollectorFunc(c.Bot, filter)
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

func (c *collectorsImpl) NewSelectMenuSubmitCollector(filter SelectMenuInteractionFilter) (<-chan *SelectMenuInteraction, func()) {
	return c.NewSelectMenuSubmitCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewSlashCommandCollector(filter SlashCommandInteractionFilter) (<-chan *SlashCommandInteraction, func()) {
	return c.NewSlashCommandCollectorFunc(c.Bot, filter)
}

func (c *collectorsImpl) NewUserCommandCollector(filter UserCommandInteractionFilter) (<-chan *UserCommandInteraction, func()) {
	return c.NewUserCommandCollectorFunc(c.Bot, filter)
}
