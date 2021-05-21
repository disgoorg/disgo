package api

import "time"

type IntegrationType string

const (
	IntegrationTypeTwitch  IntegrationType = "twitch"
	IntegrationTypeYouTube IntegrationType = "youtube"
	IntegrationTypeDiscord IntegrationType = "discord"
)

type ExpireBehaviors int

const (
	ExpireBehaviorsRemoveRole = iota
	ExpireBehaviorsKick
)

type Integration struct {
	ID                Snowflake           `json:"id"`
	Name              string              `json:"name"`
	Type              IntegrationType     `json:"type"`
	Enabled           bool                `json:"enabled"`
	Syncing           *bool               `json:"syncing"`
	RoleID            *Snowflake          `json:"role_id"`
	EnableEmoticons   *bool               `json:"enable_emoticons"`
	ExpireBehaviors   *ExpireBehaviors    `json:"expire_behaviors"`
	ExpireGracePeriod *int                `json:"expire_grace_period"`
	User              *User               `json:"user"`
	Account           *IntegrationAccount `json:"account"`
	SyncedAt          *time.Time          `json:"synced_at"`
	SubscriberCount   *int                `json:"subscriber_count"`
	Revoked           *bool               `json:"revoked"`
	Application       *Application        `json:"application"`
}

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
