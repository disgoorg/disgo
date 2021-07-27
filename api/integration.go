package api

import "github.com/DisgoOrg/restclient"

// IntegrationAccount (https://discord.com/developers/docs/resources/guild#integration-account-object)
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IntegrationApplication (https://discord.com/developers/docs/resources/guild#integration-application-object)
type IntegrationApplication struct {
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Summary     string    `json:"summary"`
	Bot         *User     `json:"bot"`
}

// Integration (https://discord.com/developers/docs/resources/guild#integration-object)
type Integration struct {
	Disgo             Disgo                   `json:"-"`
	GuildID           Snowflake               `json:"-"`
	ID                Snowflake               `json:"id"`
	Name              string                  `json:"name"`
	Type              IntegrationType         `json:"type"`
	Enabled           bool                    `json:"enabled"`
	Syncing           *bool                   `json:"syncing"`
	RoleID            *Snowflake              `json:"role_id"`
	EnableEmoticons   *bool                   `json:"enable_emoticons"`
	ExpireBehavior    *int                    `json:"expire_behavior"`
	ExpireGracePeriod *int                    `json:"expire_grace_period"`
	User              *User                   `json:"user"`
	Account           IntegrationAccount      `json:"account"`
	SyncedAt          *string                 `json:"synced_at"`
	SubscriberCount   *int                    `json:"subscriber_account"`
	Revoked           *bool                   `json:"revoked"`
	Application       *IntegrationApplication `json:"application"`
}

func (i *Integration) Guild() *Guild {
	return i.Disgo.Cache().Guild(i.GuildID)
}

func (i *Integration) Member() *Member {
	if i.User == nil {
		return nil
	}
	return i.Disgo.Cache().Member(i.GuildID, i.User.ID)
}

func (i *Integration) Role() *Role {
	if i.RoleID == nil {
		return nil
	}
	return i.Disgo.Cache().Role(*i.RoleID)
}

func (i *Integration) Delete() restclient.RestError {
	return i.Disgo.RestClient().DeleteIntegration(i.GuildID, i.ID)
}

type IntegrationType string

const (
	IntegrationTypeTwitch  IntegrationType = "twitch"
	IntegrationTypeYouTube IntegrationType = "youtube"
	IntegrationTypeDiscord IntegrationType = "discord"
)
