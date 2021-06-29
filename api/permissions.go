package api

import (
	"encoding/json"
	"strconv"
)

// PermissionOverwriteType is the type of a PermissionOverwrite
type PermissionOverwriteType int

// Constants for PermissionOverwriteType
const (
	PermissionOverwriteTypeRole PermissionOverwriteType = iota
	PermissionOverwriteTypeMember
)

// PermissionOverwrite is used to determine who can perform particular actions in a GuildChannel
type PermissionOverwrite struct {
	ID    Snowflake               `json:"id"`
	Type  PermissionOverwriteType `json:"type"`
	Allow Permissions             `json:"allow,string"`
	Deny  Permissions             `json:"deny,string"`
}

// Permissions extends the Bit structure, and is used within roles and channels
type Permissions int64

// Constants for the different bit offsets of text channel permissions
const (
	PermissionSendMessages Permissions = 1 << (iota + 11)
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
)

// Constants for the different bit offsets of voice permissions
const (
	PermissionVoiceConnect Permissions = 1 << (iota + 20)
	PermissionVoiceSpeak
	PermissionVoiceMuteMembers
	PermissionVoiceDeafenMembers
	PermissionVoiceMoveMembers
	PermissionVoiceUseVAD
	PermissionVoicePrioritySpeaker Permissions = 1 << (iota + 2)
)

// Constants for general management.
const (
	PermissionChangeNickname Permissions = 1 << (iota + 26)
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageEmojis
	PermissionUseCommands
)

// Constants for the different bit offsets of general permissions
const (
	PermissionCreateInstantInvite Permissions = 1 << iota
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageServer
	PermissionAddReactions
	PermissionViewAuditLogs
	PermissionViewChannel Permissions = 1 << (iota + 2)

	PermissionsAllText = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionsAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionsAllChannel = PermissionsAllText |
		PermissionsAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionsAll = PermissionsAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis

	PermissionsNone Permissions = 0
)

// MarshalJSON marshals permissions into a string
func (p Permissions) MarshalJSON() ([]byte, error) {
	strPermissions := strconv.FormatInt(int64(p), 10)

	jsonValue, err := json.Marshal(strPermissions)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}

// UnmarshalJSON unmarshalls permissions into an int64
func (p *Permissions) UnmarshalJSON(b []byte) error {
	var strPermissions string
	err := json.Unmarshal(b, &strPermissions)
	if err != nil {
		return err
	}

	intPermissions, err := strconv.Atoi(strPermissions)
	if err != nil {
		return err
	}
	*p = Permissions(intPermissions)
	return nil
}

// Add allows you to add multiple bits together, producing a new bit
func (p Permissions) Add(bits ...Permissions) Permissions {
	total := Permissions(0)
	for _, bit := range bits {
		total |= bit
	}
	p |= total
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Permissions) Remove(bits ...Permissions) Permissions {
	total := Permissions(0)
	for _, bit := range bits {
		total |= bit
	}
	p &^= total
	return p
}

// HasAll will ensure that the bit includes all of the bits entered
func (p Permissions) HasAll(bits ...Permissions) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (p Permissions) Has(bit Permissions) bool {
	return (p & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (p Permissions) MissingAny(bits ...Permissions) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (p Permissions) Missing(bit Permissions) bool {
	return !p.Has(bit)
}

// GetMemberPermissions returns all Permissions from the provided Member
func GetMemberPermissions(member *Member) Permissions {
	if member.IsOwner() {
		return PermissionsAll
	}
	if guild := member.Guild(); guild != nil {
		var permissions Permissions
		for _, role := range member.Roles() {
			permissions = permissions.Add(role.Permissions)
			if permissions.Has(PermissionAdministrator) {
				return PermissionsAll
			}
		}
		return permissions
	}
	return PermissionsNone
}
