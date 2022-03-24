package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

// PermissionOverwriteType is the type of PermissionOverwrite
type PermissionOverwriteType int

// Constants for PermissionOverwriteType
const (
	PermissionOverwriteTypeRole PermissionOverwriteType = iota
	PermissionOverwriteTypeMember
)

var permissionOverwrites = map[PermissionOverwriteType]func() PermissionOverwrite{
	PermissionOverwriteTypeRole: func() PermissionOverwrite {
		return &RolePermissionOverwrite{}
	},
	PermissionOverwriteTypeMember: func() PermissionOverwrite {
		return &MemberPermissionOverwrite{}
	},
}

// PermissionOverwrite is used to determine who can perform particular actions in a GetGuildChannel
type PermissionOverwrite interface {
	Type() PermissionOverwriteType
	ID() snowflake.Snowflake
}

type UnmarshalPermissionOverwrite struct {
	PermissionOverwrite
}

func (o *UnmarshalPermissionOverwrite) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type PermissionOverwriteType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	fn, ok := permissionOverwrites[oType.Type]
	if !ok {
		return fmt.Errorf("unkown permission overwrite with type %d received", oType.Type)
	}

	v := fn()
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	o.PermissionOverwrite = v
	return nil
}

type RolePermissionOverwrite struct {
	RoleID snowflake.Snowflake `json:"id"`
	Allow  Permissions         `json:"allow"`
	Deny   Permissions         `json:"deny"`
}

func (o RolePermissionOverwrite) ID() snowflake.Snowflake {
	return o.RoleID
}

func (o RolePermissionOverwrite) MarshalJSON() ([]byte, error) {
	type rolePermissionOverwrite RolePermissionOverwrite
	return json.Marshal(struct {
		Type PermissionOverwriteType
		rolePermissionOverwrite
	}{
		Type:                    o.Type(),
		rolePermissionOverwrite: rolePermissionOverwrite(o),
	})
}

func (o RolePermissionOverwrite) Type() PermissionOverwriteType {
	return PermissionOverwriteTypeRole
}

type MemberPermissionOverwrite struct {
	UserID snowflake.Snowflake `json:"id"`
	Allow  Permissions         `json:"allow"`
	Deny   Permissions         `json:"deny"`
}

func (o MemberPermissionOverwrite) ID() snowflake.Snowflake {
	return o.UserID
}

func (o MemberPermissionOverwrite) MarshalJSON() ([]byte, error) {
	type memberPermissionOverwrite MemberPermissionOverwrite
	return json.Marshal(struct {
		Type PermissionOverwriteType
		memberPermissionOverwrite
	}{
		Type:                      o.Type(),
		memberPermissionOverwrite: memberPermissionOverwrite(o),
	})
}

func (o MemberPermissionOverwrite) Type() PermissionOverwriteType {
	return PermissionOverwriteTypeMember
}

type PermissionOverwriteUpdate interface {
	Type() PermissionOverwriteType
}

type RolePermissionOverwriteUpdate struct {
	Allow Permissions `json:"allow"`
	Deny  Permissions `json:"deny"`
}

func (u RolePermissionOverwriteUpdate) MarshalJSON() ([]byte, error) {
	type rolePermissionOverwriteUpdate RolePermissionOverwriteUpdate
	return json.Marshal(struct {
		Type PermissionOverwriteType
		rolePermissionOverwriteUpdate
	}{
		Type:                          u.Type(),
		rolePermissionOverwriteUpdate: rolePermissionOverwriteUpdate(u),
	})
}

func (RolePermissionOverwriteUpdate) Type() PermissionOverwriteType {
	return PermissionOverwriteTypeRole
}

type MemberPermissionOverwriteUpdate struct {
	Allow Permissions `json:"allow"`
	Deny  Permissions `json:"deny"`
}

func (u MemberPermissionOverwriteUpdate) MarshalJSON() ([]byte, error) {
	type memberPermissionOverwriteUpdate MemberPermissionOverwriteUpdate
	return json.Marshal(struct {
		Type PermissionOverwriteType
		memberPermissionOverwriteUpdate
	}{
		Type:                            u.Type(),
		memberPermissionOverwriteUpdate: memberPermissionOverwriteUpdate(u),
	})
}

func (MemberPermissionOverwriteUpdate) Type() PermissionOverwriteType {
	return PermissionOverwriteTypeMember
}
