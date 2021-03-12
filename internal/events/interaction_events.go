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
}

type SlashCommandEvent struct {
	GenericInteractionEvent
	Name      string
	CommandID api.Snowflake
	Options   []api.OptionData
	Thread    api.CommandThread
}

func (e SlashCommandEvent) reply(message string,  ephemeral bool) *promise.Promise {
	return e.Disgo.RestClient().SendInteractionResponse(e.InteractionID, e.Token, api.InteractionResponse{})
}

