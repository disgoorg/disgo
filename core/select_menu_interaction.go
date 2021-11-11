package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// SelectMenuInteractionFilter used to filter SelectMenuInteraction(s) in a collectors.SelectMenuSubmitCollector
type SelectMenuInteractionFilter func(selectMenuInteraction *SelectMenuInteraction) bool

var _ Interaction = (*SelectMenuInteraction)(nil)
var _ ComponentInteraction = (*SelectMenuInteraction)(nil)

type SelectMenuInteraction struct {
	*InteractionFields
	Message  *Message
	CustomID discord.CustomID
	Values   []string
}

func (i *SelectMenuInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeComponent
}

func (i *SelectMenuInteraction) ComponentType() discord.ComponentType {
	return discord.ComponentTypeSelectMenu
}

func (i *SelectMenuInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, callbackType, callbackData, opts...)
}

func (i *SelectMenuInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.InteractionFields, messageCreate, opts...)
}

func (i *SelectMenuInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.InteractionFields, ephemeral, opts...)
}

func (i *SelectMenuInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	return getOriginal(i.InteractionFields, opts...)
}

func (i *SelectMenuInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateOriginal(i.InteractionFields, messageUpdate, opts...)
}

func (i *SelectMenuInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return deleteOriginal(i.InteractionFields, opts...)
}

func (i *SelectMenuInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createFollowup(i.InteractionFields, messageCreate, opts...)
}

func (i *SelectMenuInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateFollowup(i.InteractionFields, messageID, messageUpdate, opts...)
}

func (i *SelectMenuInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteFollowup(i.InteractionFields, messageID, opts...)
}

func (i *SelectMenuInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return update(i.InteractionFields, messageUpdate, opts...)
}

func (i *SelectMenuInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return deferUpdate(i.InteractionFields, opts...)
}

// UpdateSelectMenu updates the used SelectMenuComponent with a new SelectMenuComponent
func (i *SelectMenuInteraction) UpdateSelectMenu(selectMenu discord.SelectMenuComponent, opts ...rest.RequestOpt) error {
	return updateComponent(i.InteractionFields, i.Message, i.CustomID, selectMenu, opts...)
}

// SelectMenuComponent returns the SelectMenuComponent which issued this SelectMenuInteraction
func (i *SelectMenuInteraction) SelectMenuComponent() discord.SelectMenuComponent {
	// this should never be nil
	return *i.Message.SelectMenuByID(i.CustomID)
}

// SelectedOptions returns the selected SelectMenuOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []discord.SelectMenuOption {
	options := make([]discord.SelectMenuOption, len(i.Values))
	for ii, option := range i.SelectMenuComponent().Options {
		for _, value := range i.Values {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}
