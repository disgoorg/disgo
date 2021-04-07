package api

import (
	"encoding/json"
	"strconv"
)

// Permissions extends the Bit structure, and is used within roles and channels
type Permissions int64

// MarshalJSON marshals permissions into a string
func (p Permissions) MarshalJSON() ([]byte, error) {
	strPermissions := strconv.FormatInt(int64(p), 10)

	jsonValue, err := json.Marshal(strPermissions)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}

// UnmarshalJSON unmarshals permissions into an int64
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
func (p Permissions) Add(bits ...Bit) Bit {
	total := Permissions(0)
	for _, bit := range bits {
		total |= bit.(Permissions)
	}
	p |= total
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Permissions) Remove(bits ...Bit) Bit {
	total := Permissions(0)
	for _, bit := range bits {
		total |= bit.(Permissions)
	}
	p &^= total
	return p
}

// HasAll will ensure that the bit includes all of the bits entered
func (p Permissions) HasAll(bits ...Bit) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (p Permissions) Has(bit Bit) bool {
	return (p & bit.(Permissions)) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (p Permissions) MissingAny(bits ...Bit) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (p Permissions) Missing(bit Bit) bool {
	return !p.Has(bit)
}

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
	PermissionUseSlashCommands
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

	PermissionAllText = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionAllChannel = PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis
)
