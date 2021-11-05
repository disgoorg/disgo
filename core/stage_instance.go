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

func (i *StageInstance) Channel() *GuildStageVoiceChannel {
	if ch := i.Bot.Caches.ChannelCache().Get(i.ChannelID); ch != nil {
		return ch.(*GuildStageVoiceChannel)
	}
	return nil
}

func (i *StageInstance) GetSpeakers() []*Member {
	ch := i.Channel()
	if ch == nil {
		return nil
	}
	var speakers []*Member
	for _, member := range ch.Members() {
		if member.VoiceState() != nil && !member.VoiceState().Suppress {
			speakers = append(speakers)
		}
	}
	return speakers
}

func (i *StageInstance) GetListeners() []*Member {
	ch := i.Channel()
	if ch == nil {
		return nil
	}
	var listeners []*Member
	for _, member := range ch.Members() {
		if member.VoiceState() != nil && member.VoiceState().Suppress {
			listeners = append(listeners)
		}
	}
	return listeners
}

func (s *VoiceState) UpdateVoiceState(suppress *discord.OptionalBool, requestToSpeak *discord.OptionalTime, opts ...rest.RequestOpt) error {
	if s.ChannelID == nil {
		return discord.ErrMemberMustBeConnectedToChannel
	}
	return s.Bot.RestServices.GuildService().UpdateCurrentUserVoiceState(s.GuildID, discord.UserVoiceStateUpdate{ChannelID: *s.ChannelID, Suppress: suppress, RequestToSpeakTimestamp: requestToSpeak}, opts...)
}

func (i *StageInstance) Update(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, error) {
	stageInstance, err := i.Bot.RestServices.StageInstanceService().UpdateStageInstance(i.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (i *StageInstance) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.StageInstanceService().DeleteStageInstance(i.ID, opts...)
}
