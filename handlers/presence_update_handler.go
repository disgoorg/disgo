package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerPresenceUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventPresenceUpdate) {
	/*oldPresence := client.Caches().Presences().GetCopy(event.GuildID, event.PresenceUser.ID)

	_ = bot.EntityBuilder.CreatePresence(event, core.CacheStrategyYes)

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

	if oldStatus != event.Status {
		client.EventManager().DispatchEvent(&events.UserStatusUpdate{
			GenericEvent: genericEvent,
			UserID:       event.PresenceUser.ID,
			OldStatus:    oldStatus,
			Status:       event.Status,
		})
	}

	if oldClientStatus == nil || oldClientStatus.Desktop != event.ClientStatus.Desktop || oldClientStatus.Mobile != event.ClientStatus.Mobile || oldClientStatus.Web != event.ClientStatus.Web {
		client.EventManager().DispatchEvent(&events.UserClientStatusUpdate{
			GenericEvent:    genericEvent,
			UserID:          event.PresenceUser.ID,
			OldClientStatus: oldClientStatus,
			ClientStatus:    event.ClientStatus,
		})
	}

	genericUserActivity := events.GenericUserActivity{
		GenericEvent: genericEvent,
		UserID:       event.PresenceUser.ID,
		GuildID:      event.GuildID,
	}

	for _, oldActivity := range oldActivities {
		var found bool
		for _, newActivity := range event.Activities {
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

	for _, newActivity := range event.Activities {
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

	for _, newActivity := range event.Activities {
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
