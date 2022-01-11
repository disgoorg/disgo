package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgo/rest"
)

type VoiceState struct {
	discord.VoiceState
	Bot    Bot
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
func (s *VoiceState) Guild() (Guild, bool) {
	return s.Bot.Caches().Guilds().Get(s.GuildID)
}

// Channel returns the Channel of this VoiceState from the Caches
func (s *VoiceState) Channel() GuildAudioChannel {
	if s.ChannelID == nil {
		return nil
	}
	if ch := s.Bot.Caches().Channels().Get(*s.ChannelID); ch != nil {
		return ch.(GuildAudioChannel)
	}
	return nil
}

func (s *VoiceState) Update(suppress *bool, requestToSpeak *json.Nullable[discord.Time], opts ...rest.RequestOpt) error {
	if s.ChannelID == nil {
		return discord.ErrMemberMustBeConnectedToChannel
	}
	userVoiceUpdate := discord.UserVoiceStateUpdate{ChannelID: *s.ChannelID, Suppress: suppress, RequestToSpeakTimestamp: requestToSpeak}
	if s.UserID == s.Bot.ClientID {
		return s.Bot.RestServices().GuildService().UpdateCurrentUserVoiceState(s.GuildID, userVoiceUpdate, opts...)
	}
	return s.Bot.RestServices().GuildService().UpdateUserVoiceState(s.GuildID, s.UserID, userVoiceUpdate, opts...)
}
