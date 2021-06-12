package api

// VoiceServerUpdate from Discord
type VoiceServerUpdate struct {
	Token    string    `json:"token"`
	GuildID  Snowflake `json:"guild_id"`
	Endpoint *string   `json:"endpoint"`
}

// VoiceState from Discord
type VoiceState struct {
	Disgo         Disgo
	GuildID       Snowflake  `json:"guild_id"`
	ChannelID     *Snowflake `json:"channel_id"`
	UserID        Snowflake  `json:"user_id"`
	SessionID     string     `json:"session_id"`
	GuildDeafened bool       `json:"deaf"`
	GuildMuted    bool       `json:"mute"`
	SelfDeafened  bool       `json:"self_deaf"`
	SelfMuted     bool       `json:"self_mute"`
	Stream        bool       `json:"self_stream"`
	Video         bool       `json:"self_video"`
	Suppressed    bool       `json:"suppress"`
}

// Muted returns if the Member is muted
func (s *VoiceState) Muted() bool {
	return s.GuildMuted || s.SelfMuted
}

// Deafened returns if the Member is deafened
func (s *VoiceState) Deafened() bool {
	return s.GuildDeafened || s.SelfDeafened
}

// Member returns the Member of this VoiceState from the Cache
func (s *VoiceState) Member() *Member {
	return s.Disgo.Cache().Member(s.GuildID, s.UserID)
}

// User returns the User of this VoiceState from the Cache
func (s *VoiceState) User() *User {
	return s.Disgo.Cache().User(s.UserID)
}

// Guild returns the Guild of this VoiceState from the Cache
func (s *VoiceState) Guild() *Guild {
	return s.Disgo.Cache().Guild(s.GuildID)
}

// VoiceChannel returns the VoiceChannel of this VoiceState from the Cache
func (s *VoiceState) VoiceChannel() VoiceChannel {
	if s.ChannelID == nil {
		return nil
	}
	return s.Disgo.Cache().VoiceChannel(*s.ChannelID)
}
