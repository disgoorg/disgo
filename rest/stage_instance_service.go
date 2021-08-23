package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewStageInstanceService(client Client) StageInstanceService {
	return nil
}

type StageInstanceService interface {
	Service
	GetStageInstance(stageInstanceID discord.Snowflake) (*discord.StageInstance, Error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate) (*discord.StageInstance, Error)
	UpdateStageInstance(stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (*discord.StageInstance, Error)
	DeleteStageInstance(stageInstanceID discord.Snowflake) Error
}
