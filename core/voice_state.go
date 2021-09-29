package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type VoiceState struct {
	discord.VoiceState
	Bot    *Bot
	Member *Member
}

// Mute returns whether the Member is muted
func (s *VoiceState) Mute() bool {
	return s.GuildMute || s.SelfMute
}

// Deaf returns whether the Member is deafened
func (s *VoiceState) Deaf() bool {
	return s.GuildDeaf || s.SelfDeaf
}

// Guild returns the Guild of this VoiceState.
// This will only check cached guilds!
func (s *VoiceState) Guild() *Guild {
	return s.Bot.Caches.GuildCache().Get(s.GuildID)
}

// Channel returns the Channel of this VoiceState.
// This will only check cached channels!
func (s *VoiceState) Channel() *Channel {
	if s.ChannelID == nil {
		return nil
	}
	return s.Bot.Caches.ChannelCache().Get(*s.ChannelID)
}

func (s *VoiceState) Update(suppress *discord.OptionalBool, requestToSpeak *discord.OptionalTime, opts ...rest.RequestOpt) rest.Error {
	if s.ChannelID == nil {
		return rest.NewError(nil, discord.ErrMemberMustBeConnectedToChannel)
	}
	userVoiceUpdate := discord.UserVoiceStateUpdate{ChannelID: *s.ChannelID, Suppress: suppress, RequestToSpeakTimestamp: requestToSpeak}
	if s.UserID == s.Bot.ClientID {
		return s.Bot.RestServices.GuildService().UpdateCurrentUserVoiceState(s.GuildID, userVoiceUpdate, opts...)
	}
	return s.Bot.RestServices.GuildService().UpdateUserVoiceState(s.GuildID, s.UserID, userVoiceUpdate, opts...)
}
