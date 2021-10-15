package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// SelectMenuInteractionFilter used to filter SelectMenuInteraction(s) in a collectors.SelectMenuSubmitCollector
type SelectMenuInteractionFilter func(selectMenuInteraction *SelectMenuInteraction) bool

type SelectMenuInteraction struct {
	discord.SelectMenuInteraction
	FollowupInteraction
	UpdateInteraction
	Message  *Message
	CustomID string
	Values   []string
}

// UpdateSelectMenu updates the used SelectMenu with a new SelectMenu
func (i *SelectMenuInteraction) UpdateSelectMenu(selectMenu discord.SelectMenu, opts ...rest.RequestOpt) error {
	return i.UpdateComponent(selectMenu, opts...)
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
