package discord

// ApplicationCommandPermissionType is the type of the ApplicationCommandPermission
type ApplicationCommandPermissionType int

// types of ApplicationCommandPermissionType
//goland:noinspection GoUnusedConst
const (
	ApplicationCommandPermissionTypeRole = iota + 1
	ApplicationCommandPermissionTypeUser
)

// ApplicationCommandPermissions holds all permissions for a ApplicationCommand
type ApplicationCommandPermissions struct {
	ID            Snowflake                      `json:"id"`
	ApplicationID Snowflake                      `json:"application_id"`
	GuildID       Snowflake                      `json:"guild_id"`
	Permissions   []ApplicationCommandPermission `json:"permissions"`
}

// ApplicationCommandPermission holds a User or Role and if they are allowed to use the ApplicationCommand
type ApplicationCommandPermission struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

// ApplicationCommandPermissionsSet is used to bulk overwrite all ApplicationCommandPermissions
type ApplicationCommandPermissionsSet struct {
	ID          Snowflake                      `json:"id"`
	Permissions []ApplicationCommandPermission `json:"permissions"`
}
