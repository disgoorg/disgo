package rpc

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type EventDataReady struct {
	V      int          `json:"v"`
	Config ServerConfig `json:"config"`
	User   discord.User `json:"user"`
}

func (EventDataReady) messageData() {}

type ServerConfig struct {
	CDNHost     string `json:"cdn_host"`
	APIEndpoint string `json:"api_endpoint"`
	Environment string `json:"environment"`
}

type EventDataMessageCreate struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Message   discord.Message `json:"message"`
}

func (EventDataMessageCreate) messageData() {}

type EventDataMessageUpdate struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Message   discord.Message `json:"message"`
}

func (EventDataMessageUpdate) messageData() {}

type EventDataMessageDelete struct {
	ChannelID snowflake.ID `json:"channel_id"`
	Message   struct {
		ID snowflake.ID `json:"id"`
	} `json:"message"`
}

func (EventDataMessageDelete) messageData() {}

var _ error = (*EventDataError)(nil)

type EventDataError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e EventDataError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (EventDataError) messageData() {}
