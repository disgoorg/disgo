package discord

import "github.com/disgoorg/snowflake/v2"

type GuildOnboarding struct {
	GuildID           snowflake.ID            `json:"guild_id"`
	Prompts           []GuildOnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []snowflake.ID          `json:"default_channel_ids"`
	Enabled           bool                    `json:"enabled"`
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
	ID          snowflake.ID   `json:"id"`
	ChannelIDs  []snowflake.ID `json:"channel_ids"`
	RoleIDs     []snowflake.ID `json:"role_ids"`
	Emoji       PartialEmoji   `json:"emoji"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
}

type GuildOnboardingPromptType int

const (
	GuildOnboardingPromptTypeMultipleChoice GuildOnboardingPromptType = iota
	GuildOnboardingPromptTypeDropdown
)
