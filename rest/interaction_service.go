package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ InteractionService = (*InteractionServiceImpl)(nil)

func NewInteractionService(client Client) InteractionService {
	return &InteractionServiceImpl{restClient: client}
}

type InteractionService interface {
	Service
	GetInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string) (*discord.Message, Error)
	CreateInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) Error
	UpdateInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string) Error

	CreateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate) (*discord.Message, Error)
	UpdateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake) Error
}

type InteractionServiceImpl struct {
	restClient Client
}

func (s *InteractionServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *InteractionServiceImpl) GetInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.GetInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &message)
	return
}

func (s *InteractionServiceImpl) CreateInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) Error {
	compiledRoute, err := route.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return NewError(nil, err)
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return NewError(nil, err)
	}

	return s.restClient.Do(ctx, compiledRoute, body, nil)
}

func (s *InteractionServiceImpl) UpdateInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(ctx, compiledRoute, body, &message)
	return
}

func (s *InteractionServiceImpl) DeleteInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string) Error {
	compiledRoute, err := route.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(ctx, compiledRoute, nil, nil)
}

func (s *InteractionServiceImpl) CreateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(ctx, compiledRoute, body, &message)
	return
}

func (s *InteractionServiceImpl) UpdateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(ctx, compiledRoute, body, &message)
	return
}

func (s *InteractionServiceImpl) DeleteFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake) Error {
	compiledRoute, err := route.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(ctx, compiledRoute, nil, nil)
}
