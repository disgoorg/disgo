package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewInteractionService(client Client) InteractionService {
	return nil
}


type InteractionService interface {
	Service
	CreateInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) Error
	UpdateInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string) Error

	CreateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate) (*discord.Message, Error)
	UpdateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, followupMessageID discord.Snowflake) Error
}
