package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/google/go-cmp/cmp"
)

// gatewayHandlerGuildDelete handles discord.GatewayEventTypePresenceUpdate
type gatewayHandlerPresenceUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerPresenceUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypePresenceUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerPresenceUpdate) New() interface{} {
	return &discord.Presence{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerPresenceUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.Presence)

	oldPresence := bot.Caches.Presences().GetCopy(payload.GuildID, payload.PresenceUser.ID)

	_ = bot.EntityBuilder.CreatePresence(payload, core.CacheStrategyYes)

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)

	var (
		oldStatus       discord.OnlineStatus
		oldClientStatus *discord.ClientStatus
		oldActivities   []discord.Activity
	)

	if oldPresence != nil {
		oldStatus = oldPresence.Status
		oldClientStatus = &oldPresence.ClientStatus
		oldActivities = oldPresence.Activities
	}

	if oldStatus != payload.Status {
		bot.EventManager.Dispatch(&events2.UserStatusUpdateEvent{
			GenericEvent: genericEvent,
			UserID:       payload.PresenceUser.ID,
			OldStatus:    oldStatus,
			Status:       payload.Status,
		})
	}

	if oldClientStatus == nil || oldClientStatus.Desktop != payload.ClientStatus.Desktop || oldClientStatus.Mobile != payload.ClientStatus.Mobile || oldClientStatus.Web != payload.ClientStatus.Web {
		bot.EventManager.Dispatch(&events2.UserClientStatusUpdateEvent{
			GenericEvent:    genericEvent,
			UserID:          payload.PresenceUser.ID,
			OldClientStatus: oldClientStatus,
			ClientStatus:    payload.ClientStatus,
		})
	}

	genericUserActivityEvent := events2.GenericUserActivityEvent{
		GenericEvent: genericEvent,
		UserID:       payload.PresenceUser.ID,
		GuildID:      payload.GuildID,
	}

	for _, oldActivity := range oldActivities {
		var found bool
		for _, newActivity := range payload.Activities {
			if oldActivity.ID == newActivity.ID {
				found = true
				break
			}
		}
		if !found {
			genericUserActivityEvent.Activity = oldActivity
			bot.EventManager.Dispatch(&events2.UserActivityStopEvent{
				GenericUserActivityEvent: &genericUserActivityEvent,
			})
		}
	}

	for _, newActivity := range payload.Activities {
		var found bool
		for _, oldActivity := range oldActivities {
			if newActivity.ID == oldActivity.ID {
				found = true
				break
			}
		}
		if !found {
			genericUserActivityEvent.Activity = newActivity
			bot.EventManager.Dispatch(&events2.UserActivityStartEvent{
				GenericUserActivityEvent: &genericUserActivityEvent,
			})
		}
	}

	for _, newActivity := range payload.Activities {
		var oldActivity *discord.Activity
		for _, activity := range oldActivities {
			if newActivity.ID == activity.ID {
				oldActivity = &activity
				break
			}
		}
		if oldActivity != nil && !cmp.Equal(*oldActivity, newActivity) {
			genericUserActivityEvent.Activity = newActivity
			bot.EventManager.Dispatch(&events2.UserActivityUpdateEvent{
				GenericUserActivityEvent: &genericUserActivityEvent,
				OldActivity:              *oldActivity,
			})
		}
	}
}
