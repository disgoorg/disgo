package discord

import "github.com/disgoorg/snowflake"

type StagePrivacyLevel int

//goland:noinspection GoUnusedConst
const (
	StagePrivacyLevelPublic StagePrivacyLevel = iota + 1
	StagePrivacyLevelGuildOnly
)

type StageInstance struct {
	ID                   snowflake.Snowflake `json:"id"`
	GuildID              snowflake.Snowflake `json:"guild_id"`
	ChannelID            snowflake.Snowflake `json:"channel_id"`
	Topic                string              `json:"topic"`
	PrivacyLevel         StagePrivacyLevel   `json:"privacy_level"`
	DiscoverableDisabled bool                `json:"discoverable_disabled"`
}

type StageInstanceCreate struct {
	ChannelID    snowflake.Snowflake `json:"channel_id"`
	Topic        string              `json:"topic,omitempty"`
	PrivacyLevel StagePrivacyLevel   `json:"privacy_level,omitempty"`
}

type StageInstanceUpdate struct {
	Topic        *string            `json:"topic,omitempty"`
	PrivacyLevel *StagePrivacyLevel `json:"privacy_level,omitempty"`
}
