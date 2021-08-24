package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type StageInstance struct {
	discord.StageInstance
	Disgo Disgo
}

func (i *StageInstance) Guild() *Guild {
	return i.Disgo.Cache().GuildCache().Get(i.GuildID)
}

func (i *StageInstance) Channel() StageChannel {
	return i.Disgo.Cache().StageChannelCache().Get(i.ChannelID)
}

func (i *StageInstance) Update(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	stageInstance, err := i.Disgo.RestServices().StageService().UpdateStageInstance(i.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Disgo.EntityBuilder().CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (i *StageInstance) Delete(opts ...rest.RequestOpt) rest.Error {
	return i.Disgo.RestServices().StageService().DeleteStageInstance(i.ID, opts...)
}
