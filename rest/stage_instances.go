package rest

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

var _ StageInstances = (*stageInstanceImpl)(nil)

func NewStageInstances(client Client) StageInstances {
	return &stageInstanceImpl{client: client}
}

type StageInstances interface {
	GetStageInstance(channelID snowflake.ID, opts ...RequestOpt) (*discord.StageInstance, error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (*discord.StageInstance, error)
	UpdateStageInstance(channelID snowflake.ID, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (*discord.StageInstance, error)
	DeleteStageInstance(channelID snowflake.ID, opts ...RequestOpt) error
}

type stageInstanceImpl struct {
	client Client
}

func (s *stageInstanceImpl) GetStageInstance(channelID snowflake.ID, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	err = s.client.Do(GetStageInstance.Compile(nil, channelID), nil, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	err = s.client.Do(CreateStageInstance.Compile(nil), stageInstanceCreate, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) UpdateStageInstance(channelID snowflake.ID, stageInstanceUpdate discord.StageInstanceUpdate, opts ...RequestOpt) (stageInstance *discord.StageInstance, err error) {
	err = s.client.Do(UpdateStageInstance.Compile(nil, channelID), stageInstanceUpdate, &stageInstance, opts...)
	return
}

func (s *stageInstanceImpl) DeleteStageInstance(channelID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteStageInstance.Compile(nil, channelID), nil, nil, opts...)
}
