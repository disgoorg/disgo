package events

import (
	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/api"
)

type GenericInteractionEvent struct {
	api.Event
	Token         string
	InteractionID api.Snowflake
	Guild         *api.Guild
	Member        *api.Member
	User          api.User
	Channel       *api.MessageChannel
	Version       int
}

type SlashCommandEvent struct {
	GenericInteractionEvent
	Name      string
	CommandID api.Snowflake
	Options   []api.OptionData
}

func (e SlashCommandEvent) Reply(message string, ephemeral bool) *promise.Promise {
	flags := 0
	if ephemeral {
		flags = 1 << 6
	}
	return e.Disgo.RestClient().SendInteractionResponse(e.InteractionID, e.Token, api.InteractionResponse{
		Type: api.InteractionResponseTypeChannelMessageWithSource,
		Data: &api.InteractionResponseData{
			Content: message,
			Flags:   flags,
		},
	})
}
