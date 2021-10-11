package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// Interaction represents a generic Interaction received from discord
type Interaction struct {
	discord.Interaction
	Bot             *Bot
	User            *User
	Member          *Member
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
}

// Respond responds to the Interaction with the provided discord.InteractionResponse
func (i *Interaction) Respond(callbackType discord.InteractionCallbackType, data interface{}, opts ...rest.RequestOpt) error {
	response := discord.InteractionResponse{
		Type: callbackType,
		Data: data,
	}
	if i.Acknowledged {
		return discord.ErrInteractionAlreadyReplied
	}
	i.Acknowledged = true

	if !i.FromGateway() {
		i.ResponseChannel <- response
		return nil
	}

	return i.Bot.RestServices.InteractionService().CreateInteractionResponse(i.ID, i.Token, response, opts...)
}

// FromGateway returns is the Interaction came in via gateway.Gateway or httpserver.Server
func (i *Interaction) FromGateway() bool {
	return i.ResponseChannel == nil
}

// Guild returns the Guild from the Caches
func (i *Interaction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.GuildCache().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *Interaction) Channel() *Channel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Bot.Caches.ChannelCache().Get(*i.ChannelID)
}

type ApplicationCommandInteraction struct {
	*Interaction
	CreateInteractionResponses
	CommandID   discord.Snowflake
	CommandName string
	Resolved    *Resolved
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]*Channel
	Messages map[discord.Snowflake]*Message
}
