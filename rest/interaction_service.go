package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ InteractionService = (*InteractionServiceImpl)(nil)

func NewInteractionService(client Client) InteractionService {
	return &InteractionServiceImpl{restClient: client}
}

type InteractionService interface {
	Service
	GetInteractionResponse(interactionID discord.Snowflake, interactionToken string, opts ...RequestOpt) (*discord.Message, Error)
	CreateInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) Error
	UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, Error)
	DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string, opts ...RequestOpt) Error

	CreateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, Error)
	UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, Error)
	DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, opts ...RequestOpt) Error
}

type InteractionServiceImpl struct {
	restClient Client
}

func (s *InteractionServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *InteractionServiceImpl) GetInteractionResponse(interactionID discord.Snowflake, interactionToken string, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.GetInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *InteractionServiceImpl) CreateInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) Error {
	compiledRoute, err := route.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return NewError(nil, err)
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return NewError(nil, err)
	}

	return s.restClient.Do(compiledRoute, body, nil, opts...)
}

func (s *InteractionServiceImpl) UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *InteractionServiceImpl) DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *InteractionServiceImpl) CreateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *InteractionServiceImpl) UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *InteractionServiceImpl) DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}