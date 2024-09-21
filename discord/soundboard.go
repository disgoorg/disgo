package discord

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type VoiceChannelEffectAnimationType int

const (
	VoiceChannelEffectAnimationTypePremium VoiceChannelEffectAnimationType = iota
	VoiceChannelEffectAnimationTypeBasic
)

type SoundboardSound struct {
	Name      string        `json:"name"`
	SoundID   snowflake.ID  `json:"sound_id"`
	Volume    float64       `json:"volume"`
	EmojiID   *snowflake.ID `json:"emoji_id"`
	EmojiName *string       `json:"emoji_name"`
	GuildID   *snowflake.ID `json:"guild_id,omitempty"`
	Available *bool         `json:"available,omitempty"`
	User      *User         `json:"user,omitempty"`
}

func (s SoundboardSound) URL(opts ...CDNOpt) string {
	return formatAssetURL(SoundboardSoundFile, opts, s.SoundID)
}

type SoundboardSoundCreate struct {
	Name      string       `json:"name"`
	Sound     Sound        `json:"sound"`
	Volume    *float64     `json:"volume,omitempty"`
	EmojiID   snowflake.ID `json:"emoji_id,omitempty"`
	EmojiName string       `json:"emoji_name,omitempty"`
}

type SoundboardSoundUpdate struct {
	Name      *string                      `json:"name,omitempty"`
	Volume    *json.Nullable[float64]      `json:"volume,omitempty"`
	EmojiID   *json.Nullable[snowflake.ID] `json:"emoji_id,omitempty"`
	EmojiName *json.Nullable[string]       `json:"emoji_name,omitempty"`
}

type SendSoundboardSound struct {
	SoundID       snowflake.ID  `json:"sound_id"`
	SourceGuildID *snowflake.ID `json:"source_guild_id,omitempty"`
}
