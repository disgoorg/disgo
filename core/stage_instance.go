package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type StageInstance struct {
	discord.StageInstance
	Bot *Bot
}

// Guild returns the Guild this StageInstance belongs to.
// This will only check cached guilds!
func (i *StageInstance) Guild() *Guild {
	return i.Bot.Caches.Guilds().Get(i.GuildID)
}

// Channel returns the Channel this StageInstance belongs to.
func (i *StageInstance) Channel() *GuildStageVoiceChannel {
	if ch := i.Bot.Caches.Channels().Get(i.ChannelID); ch != nil {
		return ch.(*GuildStageVoiceChannel)
	}
	return nil
}

// GetSpeakers returns the Member(s) that can speak in this StageInstance
func (i *StageInstance) GetSpeakers() []*Member {
	ch := i.Channel()
	if ch == nil {
		return nil
	}
	var speakers []*Member
	for _, member := range ch.Members() {
		if member.VoiceState() != nil && !member.VoiceState().Suppress {
			speakers = append(speakers, member)
		}
	}
	return speakers
}

// GetListeners returns the Member(s) that cannot speak in this StageInstance
func (i *StageInstance) GetListeners() []*Member {
	ch := i.Channel()
	if ch == nil {
		return nil
	}
	var listeners []*Member
	for _, member := range ch.Members() {
		if member.VoiceState() != nil && member.VoiceState().Suppress {
			listeners = append(listeners, member)
		}
	}
	return listeners
}

func (s *VoiceState) UpdateVoiceState(suppress *bool, requestToSpeak *discord.NullTime, opts ...rest.RequestOpt) error {
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

// Delete deletes this StageInstance
func (i *StageInstance) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.StageInstanceService().DeleteStageInstance(i.ID, opts...)
}
