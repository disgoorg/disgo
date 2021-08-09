package discord

// CommandPermissionType is the type of the CommandPermission
type CommandPermissionType int

// types of CommandPermissionType
//goland:noinspection GoUnusedConst
const (
	CommandPermissionTypeRole = iota + 1
	CommandPermissionTypeUser
)

// GuildCommandPermissions holds all permissions for a Command
type GuildCommandPermissions struct {
	ID            Snowflake           `json:"id"`
	ApplicationID Snowflake           `json:"application_id"`
	GuildID       Snowflake           `json:"guild_id"`
	Permissions   []CommandPermission `json:"permissions"`
}

// CommandPermission holds a User or Role and if they are allowed to use the Command
type CommandPermission struct {
	ID         Snowflake             `json:"id"`
	Type       CommandPermissionType `json:"type"`
	Permission bool                  `json:"permission"`
}

// GuildCommandPermissionsSet is used to bulk overwrite all GuildCommandPermissions
type GuildCommandPermissionsSet struct {
	ID            Snowflake           `json:"id"`
	Permissions   []CommandPermission `json:"permissions"`
}
