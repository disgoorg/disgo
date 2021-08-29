package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type typingStartPayload struct {
	ChannelID discord.Snowflake
	GuildID   *discord.Snowflake
	UserID    discord.Snowflake
	Timestamp discord.Time
	Member    *discord.Member
	// TODO: check if we get user somewhere
	User discord.User
}

// TypingStartHandler handles discord.GatewayEventTypeInviteDelete
type TypingStartHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *TypingStartHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeTypingStart
}

// New constructs a new payload receiver for the raw gateway event
func (h *TypingStartHandler) New() interface{} {
	return typingStartPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *TypingStartHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(typingStartPayload)
	if !ok {
		return
	}

	user := disgo.EntityBuilder().CreateUser(payload.User, core.CacheStrategyYes)

	eventManager.Dispatch(&events.UserTypingEvent{
		GenericUserEvent: &events.GenericUserEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			UserID:       payload.UserID,
			User:         user,
		},
		ChannelID: payload.ChannelID,
	})

	if payload.GuildID == nil {
		eventManager.Dispatch(&events.DMUserTypingEvent{
			GenericUserEvent: &events.GenericUserEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				User:         user,
				UserID:       payload.UserID,
			},
			ChannelID: payload.ChannelID,
		})
	} else {
		eventManager.Dispatch(&events.GuildMemberTypingEvent{
			GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
				GenericGuildEvent: &events.GenericGuildEvent{
					GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
					GuildID:      *payload.GuildID,
					Guild:        disgo.Cache().GuildCache().Get(*payload.GuildID),
				},
				Member: disgo.EntityBuilder().CreateMember(*payload.GuildID, *payload.Member, core.CacheStrategyYes),
			},
		})
	}
}
