package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

func NewStageService(client Client) StageService {
	return &StageServiceImpl{restClient: client}
}

type StageService interface {
	Service
	GetStageInstance(channelID discord.Snowflake, opts ...RequestOpt) (*discord.StageInstance, Error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (*discord.StageInstance, Error)
	UpdateStageInstance(channelID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (*discord.StageInstance, Error)
	DeleteStageInstance(channelID discord.Snowflake, opts ...RequestOpt) Error
}

type StageServiceImpl struct {
	restClient Client
}

func (s *StageServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *StageServiceImpl) GetStageInstance(channelID discord.Snowflake, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.GetStageInstance.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &stageInstance, opts...)
	return
}

func (s *StageServiceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.CreateStageInstance.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceCreate, &stageInstance, opts...)
	return
}

func (s *StageServiceImpl) UpdateStageInstance(channelID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.UpdateStageInstance.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceUpdate, &stageInstance, opts...)
	return
}

func (s *StageServiceImpl) DeleteStageInstance(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteStageInstance.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
