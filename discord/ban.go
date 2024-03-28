package discord

import "github.com/disgoorg/snowflake/v2"

// Ban represents a banned User from a Guild (https://discord.com/developers/docs/resources/guild#ban-object)
type Ban struct {
	Reason *string `json:"reason,omitempty"`
	User   User    `json:"user"`
}

// AddBan is used to ban a User (https://discord.com/developers/docs/resources/guild#create-guild-ban-json-params)
type AddBan struct {
	DeleteMessageSeconds int `json:"delete_message_seconds,omitempty"`
}

// BulkBan is used to bulk ban Users
type BulkBan struct {
	UserIDs              []snowflake.ID `json:"user_ids"`
	DeleteMessageSeconds int            `json:"delete_message_seconds,omitempty"`
}

// BulkBanResult is the result of a BulkBan request
type BulkBanResult struct {
	BannedUsers []snowflake.ID `json:"banned_users"`
	FailedUsers []snowflake.ID `json:"failed_users"`
}
