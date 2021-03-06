package models

// Permissions extends the Bit structure, and is used within roles and channels
type Permissions Bit

// Bit returns the Bit of the Intent
func (p Permissions) Bit() Bit {
	return Bit(p)
}

// Add calls the Bit Add method
func (p Permissions) Add(bits ...Bit) Permissions {
	return Permissions(p.Bit().Add(bits...))
}

// Remove calls the Bit Remove method
func (p Permissions) Remove(bits ...Bit) Permissions {
	return Permissions(p.Bit().Remove(bits...))
}

// HasAll calls the Bit HasAll method
func (p Permissions) HasAll(bits ...Bit) bool {
	return p.Bit().HasAll(bits...)
}

// Has calls the Bit Has method
func (p Permissions) Has(bit Bit) bool {
	return p.Bit().Has(bit)
}

// MissingAny calls the Bit MissingAny method
func (p Permissions) MissingAny(bits ...Bit) bool {
	return p.Bit().MissingAny(bits...)
}

// Missing calls the Bit Missing method
func (p Permissions) Missing(bits Bit) bool {
	return p.Bit().Missing(bits)
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
