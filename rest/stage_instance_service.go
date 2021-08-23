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
	GetStageInstance(stageInstanceID discord.Snowflake) (*discord.StageInstance, Error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate) (*discord.StageInstance, Error)
	UpdateStageInstance(stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (*discord.StageInstance, Error)
	DeleteStageInstance(stageInstanceID discord.Snowflake) Error
}

type StageServiceImpl struct {
	restClient Client
}

func (s *StageServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *StageServiceImpl) GetStageInstance(channelID discord.Snowflake) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.GetStageInstance.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &stageInstance)
	return
}

func (s *StageServiceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.CreateStageInstance.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceCreate, &stageInstance)
	return
}

func (s *StageServiceImpl) UpdateStageInstance(stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (stageInstance *discord.StageInstance, rErr Error) {
	compiledRoute, err := route.UpdateStageInstance.Compile(nil, stageInstanceID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stageInstanceUpdate, &stageInstance)
	return
}

func (s *StageServiceImpl) DeleteStageInstance(stageInstanceID discord.Snowflake) Error {
	compiledRoute, err := route.DeleteStageInstance.Compile(nil, stageInstanceID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil)
}
