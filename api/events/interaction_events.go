package events

import (
	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/api"
)

type GenericInteractionEvent struct {
	api.Event
	api.Interaction
}

func (e GenericInteractionEvent) Guild() *api.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Disgo.Cache().Guild(*e.GuildID)
}

func (e GenericInteractionEvent) DMChannel() *api.DMChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().DMChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) MessageChannel() *api.MessageChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().MessageChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) TextChannel() *api.TextChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().TextChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) GuildChannel() *api.GuildChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().GuildChannel(*e.ChannelID)
}

type SlashCommandEvent struct {
	GenericInteractionEvent
	api.InteractionData
}

func (e SlashCommandEvent) Reply(message string, ephemeral bool) *promise.Promise {
	flags := 0
	if ephemeral {
		flags = 1 << 6
	}
	return e.Disgo.RestClient().SendInteractionResponse(e.Interaction.ID, e.Token, api.InteractionResponse{
		Type: api.InteractionResponseTypeChannelMessageWithSource,
		Data: &api.InteractionResponseData{
			Content: message,
			Flags:   flags,
		},
	})
}
