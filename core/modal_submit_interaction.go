package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ModalSubmitInteractionFilter func(ModalSubmitInteraction *ModalSubmitInteraction) bool

var _ Interaction = (*ModalSubmitInteraction)(nil)

type ModalSubmitInteraction struct {
	discord.ModalSubmitInteraction
	*InteractionFields
}


func (i *ModalSubmitInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, i.ID, i.Token, callbackType, callbackData, opts...)
}

func (i *ModalSubmitInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.InteractionFields, i.ID, i.Token, messageCreate, opts...)
}

func (i *ModalSubmitInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.InteractionFields, i.ID, i.Token, ephemeral, opts...)
}

func (i *ModalSubmitInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	return getOriginal(i.InteractionFields, i.ApplicationID, i.Token, opts...)
}

func (i *ModalSubmitInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateOriginal(i.InteractionFields, i.ApplicationID, i.Token, messageUpdate, opts...)
}

func (i *ModalSubmitInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return deleteOriginal(i.InteractionFields, i.ApplicationID, i.Token, opts...)
}

func (i *ModalSubmitInteraction) GetFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getFollowup(i.InteractionFields, i.ApplicationID, i.Token, messageID, opts...)
}

func (i *ModalSubmitInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createFollowup(i.InteractionFields, i.ApplicationID, i.Token, messageCreate, opts...)
}

func (i *ModalSubmitInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateFollowup(i.InteractionFields, i.ApplicationID, i.Token, messageID, messageUpdate, opts...)
}

func (i *ModalSubmitInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteFollowup(i.InteractionFields, i.ApplicationID, i.Token, messageID, opts...)
}

// Guild returns the Guild from the Caches
func (i *ModalSubmitInteraction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.Guilds().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *ModalSubmitInteraction) Channel() Channel {
	return i.Bot.Caches.Channels().Get(i.ChannelID)
}