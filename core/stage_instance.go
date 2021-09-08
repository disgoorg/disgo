package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type StageInstance struct {
	discord.StageInstance
	Bot *Bot
}

func (i *StageInstance) Guild() *Guild {
	return i.Bot.Caches.GuildCache().Get(i.GuildID)
}

func (i *StageInstance) Channel() *Channel {
	return i.Bot.Caches.ChannelCache().Get(i.ChannelID)
}

func (i *StageInstance) Update(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	stageInstance, err := i.Bot.RestServices.StageInstanceService().UpdateStageInstance(i.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (i *StageInstance) Delete(opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.StageInstanceService().DeleteStageInstance(i.ID, opts...)
}
