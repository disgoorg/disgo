package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service              = (*stageInstanceServiceImpl)(nil)
	_ StageInstanceService = (*stageInstanceServiceImpl)(nil)
)

func NewStageInstanceService(restClient Client) StageInstanceService {
	return &stageInstanceServiceImpl{restClient: restClient}
}

type StageInstanceService interface {
	Service
	GetStageInstance(guildID discord.Snowflake, opts ...RequestOpt) (*discord.StageInstance, error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (*discord.StageInstance, error)
	UpdateStageInstance(guildID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (*discord.StageInstance, error)
	DeleteStageInstance(guildID discord.Snowflake, opts ...RequestOpt) error
}

type stageInstanceServiceImpl struct {
	restClient Client
}

func (s *stageInstanceServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *stageInstanceServiceImpl) GetStageInstance(guildID discord.Snowflake, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetStageInstance.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &stageInstance, opts...)
	return
}

func (s *stageInstanceServiceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateStageInstance.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, stageInstanceCreate, &stageInstance, opts...)
	return
}

func (s *stageInstanceServiceImpl) UpdateStageInstance(guildID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateStageInstance.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, stageInstanceUpdate, &stageInstance, opts...)
	return
}

func (s *stageInstanceServiceImpl) DeleteStageInstance(guildID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteStageInstance.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
