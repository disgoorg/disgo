package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

// gatewayHandlerGuildDelete handles discord.GatewayEventTypePresenceUpdate
type gatewayHandlerPresenceUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerPresenceUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypePresenceUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerPresenceUpdate) New() any {
	return &discord.Presence{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerPresenceUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	/*payload := *v.(*discord.Presence)

	oldPresence := client.Caches().Presences().GetCopy(payload.GuildID, payload.PresenceUser.ID)

	_ = bot.EntityBuilder.CreatePresence(payload, core.CacheStrategyYes)

	genericEvent := events.NewGenericEvent(client, sequenceNumber)

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
		client.EventManager().DispatchEvent(&events.UserStatusUpdateEvent{
			GenericEvent: genericEvent,
			UserID:       payload.PresenceUser.ID,
			OldStatus:    oldStatus,
			Status:       payload.Status,
		})
	}

	if oldClientStatus == nil || oldClientStatus.Desktop != payload.ClientStatus.Desktop || oldClientStatus.Mobile != payload.ClientStatus.Mobile || oldClientStatus.Web != payload.ClientStatus.Web {
		client.EventManager().DispatchEvent(&events.UserClientStatusUpdateEvent{
			GenericEvent:    genericEvent,
			UserID:          payload.PresenceUser.ID,
			OldClientStatus: oldClientStatus,
			ClientStatus:    payload.ClientStatus,
		})
	}

	genericUserActivityEvent := events.GenericUserActivityEvent{
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
			client.EventManager().DispatchEvent(&events.UserActivityStopEvent{
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
			client.EventManager().DispatchEvent(&events.UserActivityStartEvent{
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
			client.EventManager().DispatchEvent(&events.UserActivityUpdateEvent{
				GenericUserActivityEvent: &genericUserActivityEvent,
				OldActivity:              *oldActivity,
			})
		}
	}*/
}
