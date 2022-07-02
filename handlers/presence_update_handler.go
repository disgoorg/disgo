package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerPresenceUpdate struct{}

func (h *gatewayHandlerPresenceUpdate) EventType() gateway.EventType {
	return gateway.EventTypePresenceUpdate
}

func (h *gatewayHandlerPresenceUpdate) New() any {
	return &discord.Presence{}
}

func (h *gatewayHandlerPresenceUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	/*payload := *v.(*discord.Presence)

	oldPresence := client.Caches().Presences().GetCopy(payload.GuildID, payload.PresenceUser.ID)

	_ = bot.EntityBuilder.CreatePresence(payload, core.CacheStrategyYes)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

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
		client.EventManager().DispatchEvent(&events.UserStatusUpdate{
			GenericEvent: genericEvent,
			UserID:       payload.PresenceUser.ID,
			OldStatus:    oldStatus,
			Status:       payload.Status,
		})
	}

	if oldClientStatus == nil || oldClientStatus.Desktop != payload.ClientStatus.Desktop || oldClientStatus.Mobile != payload.ClientStatus.Mobile || oldClientStatus.Web != payload.ClientStatus.Web {
		client.EventManager().DispatchEvent(&events.UserClientStatusUpdate{
			GenericEvent:    genericEvent,
			UserID:          payload.PresenceUser.ID,
			OldClientStatus: oldClientStatus,
			ClientStatus:    payload.ClientStatus,
		})
	}

	genericUserActivity := events.GenericUserActivity{
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
			client.EventManager().DispatchEvent(&events.UserActivityStop{
				GenericUserActivity: &genericUserActivityEvent,
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
			client.EventManager().DispatchEvent(&events.UserActivityStart{
				GenericUserActivity: &genericUserActivityEvent,
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
			client.EventManager().DispatchEvent(&events.UserActivityUpdate{
				GenericUserActivity: &genericUserActivityEvent,
				OldActivity:              *oldActivity,
			})
		}
	}*/
}
