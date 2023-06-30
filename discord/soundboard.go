package discord

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type SoundboardEffectAnimationType int

const (
	SoundboardEffectAnimationTypePremium SoundboardEffectAnimationType = iota
	SoundboardEffectAnimationTypeBasic
)

type SoundboardSound struct {
	Name         string        `json:"name"`
	SoundID      snowflake.ID  `json:"sound_id"`
	ID           *snowflake.ID `json:"id"`
	Volume       float64       `json:"volume"`
	EmojiID      *snowflake.ID `json:"emoji_id"`
	EmojiName    *string       `json:"emoji_name"`
	OverridePath *string       `json:"override_path"`
	GuildID      *snowflake.ID `json:"guild_id,omitempty"`
	UserID       snowflake.ID  `json:"user_id"`
	Available    *bool         `json:"available,omitempty"`
	User         *User         `json:"user"`
}

// FileURL returns the URL for this sound
func (s SoundboardSound) FileURL() string {
	if s.ID != nil {
		return SoundboardSoundURL(*s.ID)
	}
	return SoundboardDefaultSoundURL(*s.OverridePath)
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