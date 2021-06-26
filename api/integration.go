package api

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
	ID                Snowflake               `json:"id"`
	Name              string                  `json:"name"`
	Type              string                  `json:"type"`
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
