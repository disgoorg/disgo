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
	GetStageInstance(ctx context.Context, stageInstanceID discord.Snowflake) (*discord.StageInstance, Error)
	CreateStageInstance(ctx context.Context, stageInstanceCreate discord.StageInstanceCreate) (*discord.StageInstance, Error)
	UpdateStageInstance(ctx context.Context, stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (*discord.StageInstance, Error)
	DeleteStageInstance(ctx context.Context, stageInstanceID discord.Snowflake) Error
}
