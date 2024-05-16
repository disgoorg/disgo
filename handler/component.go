package handler

import (
	"context"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

type ComponentEvent struct {
	*events.ComponentInteractionCreate
	Vars map[string]string
	Ctx  context.Context
}

func (e *ComponentEvent) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ComponentEvent) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *ComponentEvent) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ComponentEvent) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *ComponentEvent) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *ComponentEvent) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *ComponentEvent) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}
