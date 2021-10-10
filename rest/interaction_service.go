package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ InteractionService = (*interactionServiceImpl)(nil)

func NewInteractionService(restClient Client) InteractionService {
	return &interactionServiceImpl{restClient: restClient}
}

type InteractionService interface {
	Service
	GetInteractionResponse(interactionID discord.Snowflake, interactionToken string, opts ...RequestOpt) (*discord.Message, error)
	CreateInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) error
	UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string, opts ...RequestOpt) error

	CreateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, error)
	UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, opts ...RequestOpt) error
}

type interactionServiceImpl struct {
	restClient Client
}

func (s *interactionServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *interactionServiceImpl) GetInteractionResponse(interactionID discord.Snowflake, interactionToken string, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *interactionServiceImpl) CreateInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) error {
	compiledRoute, err := route.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return err
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return err
	}

	return s.restClient.Do(compiledRoute, body, nil, opts...)
}

func (s *interactionServiceImpl) UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionServiceImpl) DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *interactionServiceImpl) CreateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionServiceImpl) UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionServiceImpl) DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
