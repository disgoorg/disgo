package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ StageInstances = (*stageInstanceImpl)(nil)

func NewStageInstances(client Client) StageInstances {
	return &stageInstanceImpl{client: client}
}

type StageInstances interface {
	GetStageInstance(guildID snowflake.ID, opts ...RequestOpt) (*discord.StageInstance, error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (*discord.StageInstance, error)
	UpdateStageInstance(guildID snowflake.ID, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (*discord.StageInstance, error)
	DeleteStageInstance(guildID snowflake.ID, opts ...RequestOpt) error
}

type stageInstanceImpl struct {
	client Client
}

func (s *stageInstanceImpl) GetStageInstance(guildID snowflake.ID, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetStageInstance.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateStageInstance.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, stageInstanceCreate, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) UpdateStageInstance(guildID snowflake.ID, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateStageInstance.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, stageInstanceUpdate, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) DeleteStageInstance(guildID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteStageInstance.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}
