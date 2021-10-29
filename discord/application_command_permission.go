package discord

import "github.com/DisgoOrg/disgo/json"

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
type ApplicationCommandPermission interface {
	json.Marshaler
	Type() ApplicationCommandPermissionType
}

type ApplicationCommandPermissionUser struct {
	ID         Snowflake `json:"id"`
	Permission bool      `json:"permission"`
}

func (p ApplicationCommandPermissionUser) MarshalJSON() ([]byte, error) {
	type applicationCommandPermissionUser ApplicationCommandPermissionUser
	v := struct {
		Type ApplicationCommandPermissionType `json:"type"`
		applicationCommandPermissionUser
	}{
		Type:                             p.Type(),
		applicationCommandPermissionUser: applicationCommandPermissionUser(p),
	}
	return json.Marshal(v)
}

func (p ApplicationCommandPermissionUser) Type() ApplicationCommandPermissionType {
	return ApplicationCommandPermissionTypeUser
}

type ApplicationCommandPermissionRole struct {
	ID         Snowflake `json:"id"`
	Permission bool      `json:"permission"`
}

func (p ApplicationCommandPermissionRole) MarshalJSON() ([]byte, error) {
	type applicationCommandPermissionRole ApplicationCommandPermissionRole
	v := struct {
		Type ApplicationCommandPermissionType `json:"type"`
		applicationCommandPermissionRole
	}{
		Type:                             p.Type(),
		applicationCommandPermissionRole: applicationCommandPermissionRole(p),
	}
	return json.Marshal(v)
}

func (p ApplicationCommandPermissionRole) Type() ApplicationCommandPermissionType {
	return ApplicationCommandPermissionTypeRole
}

// ApplicationCommandPermissionsSet is used to bulk overwrite all ApplicationCommandPermissions
type ApplicationCommandPermissionsSet struct {
	ID          Snowflake                      `json:"id"`
	Permissions []ApplicationCommandPermission `json:"permissions"`
}
