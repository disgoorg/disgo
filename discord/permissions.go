package discord

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/disgoorg/json"
)

var EmptyStringBytes = []byte(`""`)

// Permissions extends the Bit structure, and is used within roles and channels (https://discord.com/developers/docs/topics/permissions#permissions)
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
	PermissionManageEmojisAndStickers
	PermissionUseApplicationCommands
	PermissionRequestToSpeak
	PermissionManageEvents
	PermissionManageThreads
	PermissionCreatePublicThread
	PermissionCreatePrivateThread
	PermissionUseExternalStickers
	PermissionSendMessagesInThreads
	PermissionStartEmbeddedActivities
	PermissionModerateMembers
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

	PermissionsAllThread = PermissionManageThreads |
		PermissionCreatePublicThread |
		PermissionCreatePrivateThread |
		PermissionSendMessagesInThreads

	PermissionsAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker

	PermissionsAllChannel = PermissionsAllText |
		PermissionsAllThread |
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
		PermissionManageEmojisAndStickers

	PermissionsStageModerator = PermissionManageChannels |
		PermissionVoiceMuteMembers |
		PermissionVoiceMoveMembers

	PermissionsNone Permissions = 0
)

var permissions = map[Permissions]string{
	PermissionCreateInstantInvite:     "Create Instant Invite",
	PermissionKickMembers:             "Kick Members",
	PermissionBanMembers:              "Ban Members",
	PermissionAdministrator:           "Administrator",
	PermissionManageChannels:          "Manage Channels",
	PermissionManageServer:            "Manage Server",
	PermissionAddReactions:            "Add Reactions",
	PermissionViewAuditLogs:           "View Audit Logs",
	PermissionViewChannel:             "View Channel",
	PermissionSendMessages:            "Send Messages",
	PermissionSendTTSMessages:         "Send TTS Messages",
	PermissionManageMessages:          "Manage Messages",
	PermissionEmbedLinks:              "Embed Links",
	PermissionAttachFiles:             "Attach Files",
	PermissionReadMessageHistory:      "Read Message History",
	PermissionMentionEveryone:         "Mention Everyone",
	PermissionUseExternalEmojis:       "Use External Emojis",
	PermissionVoiceConnect:            "Connect",
	PermissionVoiceSpeak:              "Speak",
	PermissionVoiceMuteMembers:        "Mute Members",
	PermissionVoiceDeafenMembers:      "Deafen Members",
	PermissionVoiceMoveMembers:        "Move Members",
	PermissionVoiceUseVAD:             "Use Voice Activity",
	PermissionVoicePrioritySpeaker:    "Priority Speaker",
	PermissionChangeNickname:          "Change Nickname",
	PermissionManageNicknames:         "Manage Nicknames",
	PermissionManageRoles:             "Manage Roles",
	PermissionManageWebhooks:          "Manage Webhooks",
	PermissionManageEmojisAndStickers: "Manage Emojis and Stickers",
	PermissionUseApplicationCommands:  "Use Application Commands",
	PermissionRequestToSpeak:          "Request to Speak",
	PermissionManageEvents:            "Manage Events",
	PermissionManageThreads:           "Manage Threads",
	PermissionCreatePublicThread:      "Create Public Threads",
	PermissionCreatePrivateThread:     "Create Private Threads",
	PermissionUseExternalStickers:     "Use External Stickers",
	PermissionSendMessagesInThreads:   "Send Messages in Threads",
	PermissionStartEmbeddedActivities: "Start Embedded Activities",
	PermissionModerateMembers:         "Moderate Members",
}

func (p Permissions) String() string {
	if p == PermissionsNone {
		return "None"
	}
	perms := new(strings.Builder)
	for permission, name := range permissions {
		if p.Has(permission) {
			perms.WriteString(name)
			perms.WriteString(", ")
		}
	}
	return perms.String()[:perms.Len()-2] // remove trailing comma and space
}

// MarshalJSON marshals permissions into a string
func (p Permissions) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(p), 10))
}

// UnmarshalJSON unmarshalls permissions into an int64
func (p *Permissions) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, EmptyStringBytes) || bytes.Equal(data, json.NullBytes) {
		return nil
	}

	str, _ := strconv.Unquote(string(data))
	perms, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	*p = Permissions(perms)
	return nil
}

// Add allows you to add multiple bits together, producing a new bit
func (p Permissions) Add(bits ...Permissions) Permissions {
	for _, bit := range bits {
		p |= bit
	}
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Permissions) Remove(bits ...Permissions) Permissions {
	for _, bit := range bits {
		p &^= bit
	}
	return p
}

// Has will ensure that the bit includes all the bits entered
func (p Permissions) Has(bits ...Permissions) bool {
	for _, bit := range bits {
		if (p & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (p Permissions) Missing(bits ...Permissions) bool {
	for _, bit := range bits {
		if (p & bit) != bit {
			return true
		}
	}
	return false
}
