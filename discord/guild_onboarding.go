package discord

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type GuildOnboarding struct {
	GuildID           snowflake.ID            `json:"guild_id"`
	Prompts           []GuildOnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []snowflake.ID          `json:"default_channel_ids"`
	Enabled           bool                    `json:"enabled"`
	Mode              GuildOnboardingMode     `json:"mode"`
}

type GuildOnboardingPrompt struct {
	ID           snowflake.ID                  `json:"id"`
	Options      []GuildOnboardingPromptOption `json:"options"`
	Title        string                        `json:"title"`
	SingleSelect bool                          `json:"single_select"`
	Required     bool                          `json:"required"`
	InOnboarding bool                          `json:"in_onboarding"`
	Type         GuildOnboardingPromptType     `json:"type"`
}

type GuildOnboardingPromptOption struct {
	ID         snowflake.ID   `json:"id"`
	ChannelIDs []snowflake.ID `json:"channel_ids"`
	RoleIDs    []snowflake.ID `json:"role_ids"`
	// When creating or updating prompts and their options, this field will be broken down into 3 separate fields in the payload: https://github.com/discord/discord-api-docs/pull/6479
	Emoji       PartialEmoji `json:"emoji"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
}

func (o GuildOnboardingPromptOption) MarshalJSON() ([]byte, error) {
	type onboardingPromptOption GuildOnboardingPromptOption
	return json.Marshal(struct {
		EmojiID       *snowflake.ID `json:"emoji_id,omitempty"`
		EmojiName     *string       `json:"emoji_name,omitempty"`
		EmojiAnimated bool          `json:"emoji_animated"`
		onboardingPromptOption
	}{
		EmojiID:                o.Emoji.ID,
		EmojiName:              o.Emoji.Name,
		EmojiAnimated:          o.Emoji.Animated,
		onboardingPromptOption: onboardingPromptOption(o),
	})
}

type GuildOnboardingPromptType int

const (
	GuildOnboardingPromptTypeMultipleChoice GuildOnboardingPromptType = iota
	GuildOnboardingPromptTypeDropdown
)

type GuildOnboardingMode int

const (
	GuildOnboardingModeDefault GuildOnboardingMode = iota
	GuildOnboardingModeAdvanced
)

type GuildOnboardingUpdate struct {
	Prompts           *[]GuildOnboardingPrompt `json:"prompts,omitempty"`
	DefaultChannelIDs *[]snowflake.ID          `json:"default_channel_ids,omitempty"`
	Enabled           *bool                    `json:"enabled,omitempty"`
	Mode              *GuildOnboardingMode     `json:"mode,omitempty"`
}
