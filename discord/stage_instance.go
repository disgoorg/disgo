package discord

type StagePrivacyLevel int

//goland:noinspection GoUnusedConst
const (
	StagePrivacyLevelPublic StagePrivacyLevel = iota + 1
	StagePrivacyLevelGuildOnly
)

type StageInstance struct {
	ID                   Snowflake         `json:"id"`
	GuildID              Snowflake         `json:"guild_id"`
	ChannelID            Snowflake         `json:"channel_id"`
	Topic                string            `json:"topic"`
	PrivacyLevel         StagePrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled bool              `json:"discoverable_disabled"`
}

type StageInstanceCreate struct {
	ChannelID    Snowflake         `json:"channel_id"`
	Topic        string            `json:"topic,omitempty"`
	PrivacyLevel StagePrivacyLevel `json:"privacy_level,omitempty"`
}

type StageInstanceUpdate struct {
	Topic        *string            `json:"topic,omitempty"`
	PrivacyLevel *StagePrivacyLevel `json:"privacy_level,omitempty"`
}
