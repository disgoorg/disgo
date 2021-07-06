package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type webhooksUpdateData struct {
	GuildID   api.Snowflake `json:"guild_id"`
	ChannelID api.Snowflake `json:"channel_id"`
}

// WebhooksUpdateHandler handles api.GatewayEventWebhooksUpdate
type WebhooksUpdateHandler struct{}

// Event returns the raw api.GatewayEventType
func (h *WebhooksUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventWebhooksUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *WebhooksUpdateHandler) New() interface{} {
	return &webhooksUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *WebhooksUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	webhooksUpdateData, ok := i.(*webhooksUpdateData)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.WebhooksUpdateEvent{
		GenericTextChannelEvent: &events.GenericTextChannelEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
					ChannelID:    webhooksUpdateData.ChannelID,
					Channel:      disgo.Cache().Channel(webhooksUpdateData.ChannelID),
				},
				GuildID:      webhooksUpdateData.GuildID,
				GuildChannel: disgo.Cache().GuildChannel(webhooksUpdateData.ChannelID),
			},
			TextChannel: disgo.Cache().TextChannel(webhooksUpdateData.ChannelID),
		},
	})
}
