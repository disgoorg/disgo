package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Interactions = (*interactionImpl)(nil)

func NewInteractions(client Client) Interactions {
	return &interactionImpl{client: client}
}

type Interactions interface {
	GetInteractionResponse(applicationID snowflake.ID, interactionToken string, opts ...RequestOpt) (*discord.Message, error)
	CreateInteractionResponse(interactionID snowflake.ID, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) error
	UpdateInteractionResponse(applicationID snowflake.ID, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteInteractionResponse(applicationID snowflake.ID, interactionToken string, opts ...RequestOpt) error

	GetFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, opts ...RequestOpt) (*discord.Message, error)
	CreateFollowupMessage(applicationID snowflake.ID, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, error)
	UpdateFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, opts ...RequestOpt) error
}

type interactionImpl struct {
	client Client
}

func (s *interactionImpl) GetInteractionResponse(interactionID snowflake.ID, interactionToken string, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *interactionImpl) CreateInteractionResponse(interactionID snowflake.ID, interactionToken string, interactionResponse discord.InteractionResponse, opts ...RequestOpt) error {
	compiledRoute, err := route.CreateInteractionResponse.Compile(nil, interactionID, interactionToken)
	if err != nil {
		return err
	}

	body, err := interactionResponse.ToBody()
	if err != nil {
		return err
	}

	return s.client.Do(compiledRoute, body, nil, opts...)
}

func (s *interactionImpl) UpdateInteractionResponse(applicationID snowflake.ID, interactionToken string, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionImpl) DeleteInteractionResponse(applicationID snowflake.ID, interactionToken string, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteInteractionResponse.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *interactionImpl) GetFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *interactionImpl) CreateFollowupMessage(applicationID snowflake.ID, interactionToken string, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateFollowupMessage.Compile(nil, applicationID, interactionToken)
	if err != nil {
		return
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionImpl) UpdateFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *interactionImpl) DeleteFollowupMessage(applicationID snowflake.ID, interactionToken string, messageID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteFollowupMessage.Compile(nil, applicationID, interactionToken, messageID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}
