package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ StageInstanceService = (*StageInstanceServiceImpl)(nil)

func NewStageInstanceService(restClient Client) StageInstanceService {
	return &StageInstanceServiceImpl{restClient: restClient}
}

type StageInstanceService interface {
	Service
	GetStageInstance(channelID discord.Snowflake, opts ...RequestOpt) (*discord.StageInstance, Error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (*discord.StageInstance, Error)
	UpdateStageInstance(channelID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (*discord.StageInstance, Error)
	DeleteStageInstance(channelID discord.Snowflake, opts ...RequestOpt) Error
}

type StageInstanceServiceImpl struct {
	restClient Client
}

func (s *StageInstanceServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *StageInstanceServiceImpl) GetStageInstance(channelID discord.Snowflake, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.GetStageInstance.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &stageInstance, opts...)
	return
}

func (s *StageInstanceServiceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.CreateStageInstance.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceCreate, &stageInstance, opts...)
	return
}

func (s *StageInstanceServiceImpl) UpdateStageInstance(channelID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.UpdateStageInstance.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceUpdate, &stageInstance, opts...)
	return
}

func (s *StageInstanceServiceImpl) DeleteStageInstance(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteStageInstance.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
