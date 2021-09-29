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
	return i.Bot.Caches.GuildCache().Get(i.GuildID)
}

// Channel returns the Channel this StageInstance belongs to.
// This will only check cached channels!
func (i *StageInstance) Channel() *Channel {
	return i.Bot.Caches.ChannelCache().Get(i.ChannelID)
}

// GetSpeakers returns the Member(s) that can speak in this StageInstance
func (i *StageInstance) GetSpeakers() []*Member {
	var speakers []*Member
	for _, member := range i.Channel().Members() {
		if member.VoiceState() != nil && !member.VoiceState().Suppress {
			speakers = append(speakers)
		}
	}
	return speakers
}

// GetListeners returns the Member(s) that cannot speak in this StageInstance
func (i *StageInstance) GetListeners() []*Member {
	var listeners []*Member
	for _, member := range i.Channel().Members() {
		if member.VoiceState() != nil && member.VoiceState().Suppress {
			listeners = append(listeners)
		}
	}
	return listeners
}

func (s *VoiceState) UpdateVoiceState(suppress *discord.OptionalBool, requestToSpeak *discord.OptionalTime, opts ...rest.RequestOpt) rest.Error {
	if s.ChannelID == nil {
		return rest.NewError(nil, discord.ErrMemberMustBeConnectedToChannel)
	}
	return s.Bot.RestServices.GuildService().UpdateCurrentUserVoiceState(s.GuildID, discord.UserVoiceStateUpdate{ChannelID: *s.ChannelID, Suppress: suppress, RequestToSpeakTimestamp: requestToSpeak}, opts...)
}

func (i *StageInstance) Update(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	stageInstance, err := i.Bot.RestServices.StageInstanceService().UpdateStageInstance(i.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

// Delete deletes this StageInstance
func (i *StageInstance) Delete(opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.StageInstanceService().DeleteStageInstance(i.ID, opts...)
}
