package handlers

import "github.com/DisgoOrg/disgo/api"

type threadListSyncPayload struct {
	GuildID    api.Snowflake      `json:"guild_id"`
	ChannelIDs []api.Snowflake    `json:"channel_ids"`
	Threads    []*api.ChannelImpl `json:"threads"`
	Members    []*api.Member      `json:"members"`
}

type ThreadListSyncHandler struct{}

func (h ThreadListSyncHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadCreate
}

func (h ThreadListSyncHandler) New() interface{} {
	return &threadListSyncPayload{}
}

func (h ThreadListSyncHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	threadListSyncPayload, ok := i.(*threadListSyncPayload)
	if !ok {
		return
	}
}
