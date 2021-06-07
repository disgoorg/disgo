package api

// GuildCommandPermissions holds all permissions for a Command
type GuildCommandPermissions struct {
	Disgo         Disgo
	ID            Snowflake           `json:"id"`
	ApplicationID Snowflake           `json:"application_id"`
	GuildID       Snowflake           `json:"guild_id"`
	Permissions   []CommandPermission `json:"permissions"`
}

// TODO: add methods to update those

// CommandPermissionType is the type of the CommandPermission
type CommandPermissionType int

// types of CommandPermissionType
const (
	CommandPermissionTypeRole = iota + 1
	CommandPermissionTypeUser
)

// CommandPermission holds a User or Role and if they are allowed to use the Command
type CommandPermission struct {
	ID         Snowflake             `json:"id"`
	Type       CommandPermissionType `json:"type"`
	Permission bool                  `json:"permission"`
}

// SetGuildCommandsPermissions holds a slice of SetGuildCommandPermissions
type SetGuildCommandsPermissions []SetGuildCommandPermissions

// SetGuildCommandPermissions is used to update CommandPermission ID should be omitted fro bulk update
type SetGuildCommandPermissions struct {
	ID          Snowflake           `json:"id,omitempty"`
	Permissions []CommandPermission `json:"permissions"`
}
