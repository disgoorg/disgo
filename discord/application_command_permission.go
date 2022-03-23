package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

// ApplicationCommandPermissionType is the type of the ApplicationCommandPermission
type ApplicationCommandPermissionType int

// types of ApplicationCommandPermissionType
//goland:noinspection GoUnusedConst
const (
	ApplicationCommandPermissionTypeRole = iota + 1
	ApplicationCommandPermissionTypeUser
)

// ApplicationCommandPermissionsSet is used to bulk overwrite all ApplicationCommandPermissions
type ApplicationCommandPermissionsSet struct {
	ID          snowflake.Snowflake            `json:"id,omitempty"`
	Permissions []ApplicationCommandPermission `json:"permissions"`
}

// ApplicationCommandPermissions holds all permissions for a ApplicationCommand
type ApplicationCommandPermissions struct {
	ID            snowflake.Snowflake            `json:"id"`
	ApplicationID snowflake.Snowflake            `json:"application_id"`
	GuildID       snowflake.Snowflake            `json:"guild_id"`
	Permissions   []ApplicationCommandPermission `json:"permissions"`
}

func (p *ApplicationCommandPermissions) UnmarshalJSON(data []byte) error {
	type applicationCommandPermissions ApplicationCommandPermissions
	var v struct {
		Permissions []UnmarshalApplicationCommandPermission `json:"permissions"`
		applicationCommandPermissions
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*p = ApplicationCommandPermissions(v.applicationCommandPermissions)

	if len(v.Permissions) > 0 {
		p.Permissions = make([]ApplicationCommandPermission, len(v.Permissions))
		for i := range v.Permissions {
			p.Permissions[i] = v.Permissions[i].ApplicationCommandPermission
		}
	}
	return nil
}

type UnmarshalApplicationCommandPermission struct {
	ApplicationCommandPermission
}

func (p *UnmarshalApplicationCommandPermission) UnmarshalJSON(data []byte) error {
	var pType struct {
		Type ApplicationCommandPermissionType `json:"type"`
	}

	if err := json.Unmarshal(data, &pType); err != nil {
		return err
	}

	var (
		applicationCommandPermission ApplicationCommandPermission
		err                          error
	)

	switch pType.Type {
	case ApplicationCommandPermissionTypeRole:
		var v ApplicationCommandPermissionRole
		err = json.Unmarshal(data, &v)
		applicationCommandPermission = v

	case ApplicationCommandPermissionTypeUser:
		var v ApplicationCommandPermissionUser
		err = json.Unmarshal(data, &v)
		applicationCommandPermission = v

	default:
		err = fmt.Errorf("unkown application command permission with type %d received", pType.Type)
	}

	if err != nil {
		return err
	}

	p.ApplicationCommandPermission = applicationCommandPermission
	return nil
}

// ApplicationCommandPermission holds a User or Role and if they are allowed to use the ApplicationCommand
type ApplicationCommandPermission interface {
	json.Marshaler
	Type() ApplicationCommandPermissionType
	ID() snowflake.Snowflake
	applicationCommandPermission()
}

type ApplicationCommandPermissionUser struct {
	UserID     snowflake.Snowflake `json:"id"`
	Permission bool                `json:"permission"`
}

func (p ApplicationCommandPermissionUser) MarshalJSON() ([]byte, error) {
	type applicationCommandPermissionUser ApplicationCommandPermissionUser
	return json.Marshal(struct {
		Type ApplicationCommandPermissionType `json:"type"`
		applicationCommandPermissionUser
	}{
		Type:                             p.Type(),
		applicationCommandPermissionUser: applicationCommandPermissionUser(p),
	})
}

func (ApplicationCommandPermissionUser) Type() ApplicationCommandPermissionType {
	return ApplicationCommandPermissionTypeUser
}

func (p ApplicationCommandPermissionUser) ID() snowflake.Snowflake {
	return p.UserID
}

func (ApplicationCommandPermissionUser) applicationCommandPermission() {}

type ApplicationCommandPermissionRole struct {
	RoleID     snowflake.Snowflake `json:"id"`
	Permission bool                `json:"permission"`
}

func (p ApplicationCommandPermissionRole) MarshalJSON() ([]byte, error) {
	type applicationCommandPermissionRole ApplicationCommandPermissionRole
	return json.Marshal(struct {
		Type ApplicationCommandPermissionType `json:"type"`
		applicationCommandPermissionRole
	}{
		Type:                             p.Type(),
		applicationCommandPermissionRole: applicationCommandPermissionRole(p),
	})
}

func (ApplicationCommandPermissionRole) Type() ApplicationCommandPermissionType {
	return ApplicationCommandPermissionTypeRole
}

func (p ApplicationCommandPermissionRole) ID() snowflake.Snowflake {
	return p.RoleID
}

func (ApplicationCommandPermissionRole) applicationCommandPermission() {}
