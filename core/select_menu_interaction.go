package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// SelectMenuInteractionFilter used to filter SelectMenuInteraction(s) in a collectors.SelectMenuSubmitCollector
type SelectMenuInteractionFilter func(selectMenuInteraction *SelectMenuInteraction) bool

type SelectMenuInteraction struct {
	discord.SelectMenuInteraction
	Bot             *Bot
	User            *User
	Member          *Member
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
	Message         *Message
	CustomID        string
	Values          []string
}

func (i *SelectMenuInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, callbackType, callbackData, opts...)
}

func (i *SelectMenuInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, messageCreate, opts...)
}

func (i *SelectMenuInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, ephemeral, opts...)
}

func (i *SelectMenuInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return update(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, messageUpdate, opts...)
}

func (i *SelectMenuInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return deferUpdate(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, opts...)
}

// UpdateSelectMenu updates the used SelectMenu with a new SelectMenu
func (i *SelectMenuInteraction) UpdateSelectMenu(selectMenu discord.SelectMenu, opts ...rest.RequestOpt) error {
	return updateComponent(i.Bot, i.ID, i.Token, i.ResponseChannel, &i.Acknowledged, i.Message, i.CustomID, selectMenu, opts...)
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction
func (i *SelectMenuInteraction) SelectMenu() discord.SelectMenu {
	// this should never be nil
	return *i.Message.SelectMenuByID(i.CustomID)
}

// SelectedOptions returns the selected SelectMenuOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []discord.SelectMenuOption {
	options := make([]discord.SelectMenuOption, len(i.Values))
	for ii, option := range i.SelectMenu().Options {
		for _, value := range i.Values {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}
