package discord

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/disgoorg/json"

	"github.com/disgoorg/disgo/internal/flags"
)

var EmptyStringBytes = []byte(`""`)

// Permissions extends the Bit structure, and is used within roles and channels (https://discord.com/developers/docs/topics/permissions#permissions)
type Permissions int64

const (
	PermissionCreateInstantInvite Permissions = 1 << iota
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageGuild
	PermissionAddReactions
	PermissionViewAuditLog
	PermissionPrioritySpeaker
	PermissionStream
	PermissionViewChannel
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
	PermissionViewGuildInsights
	PermissionConnect
	PermissionSpeak
	PermissionMuteMembers
	PermissionDeafenMembers
	PermissionMoveMembers
	PermissionUseVAD
	PermissionChangeNickname
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageGuildExpressions
	PermissionUseApplicationCommands
	PermissionRequestToSpeak
	PermissionManageEvents
	PermissionManageThreads
	PermissionCreatePublicThreads
	PermissionCreatePrivateThreads
	PermissionUseExternalStickers
	PermissionSendMessagesInThreads
	PermissionUseEmbeddedActivities
	PermissionModerateMembers
	PermissionViewCreatorMonetizationAnalytics
	PermissionUseSoundboard
	PermissionCreateGuildExpressions
	PermissionCreateEvents
	PermissionUseExternalSounds
	PermissionSendVoiceMessages
	_
	_
	PermissionSendPolls
	PermissionUseExternalApps

	PermissionsAllText = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone |
		PermissionSendVoiceMessages |
		PermissionSendPolls |
		PermissionUseExternalApps

	PermissionsAllThread = PermissionManageThreads |
		PermissionCreatePublicThreads |
		PermissionCreatePrivateThreads |
		PermissionSendMessagesInThreads

	PermissionsAllVoice = PermissionViewChannel |
		PermissionConnect |
		PermissionSpeak |
		PermissionStream |
		PermissionMuteMembers |
		PermissionDeafenMembers |
		PermissionMoveMembers |
		PermissionUseVAD |
		PermissionPrioritySpeaker |
		PermissionUseSoundboard |
		PermissionUseExternalSounds |
		PermissionRequestToSpeak |
		PermissionUseEmbeddedActivities |
		PermissionCreateGuildExpressions |
		PermissionCreateEvents |
		PermissionManageEvents

	PermissionsAllChannel = PermissionsAllText |
		PermissionsAllThread |
		PermissionsAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionUseExternalEmojis |
		PermissionUseApplicationCommands |
		PermissionUseExternalStickers

	PermissionsAll = PermissionsAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageGuild |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageGuildExpressions |
		PermissionViewCreatorMonetizationAnalytics |
		PermissionViewGuildInsights |
		PermissionViewAuditLog |
		PermissionManageRoles |
		PermissionChangeNickname |
		PermissionManageNicknames |
		PermissionModerateMembers

	PermissionsNone Permissions = 0
)

var permissions = map[Permissions]string{
	PermissionCreateInstantInvite:              "Create Instant Invite",
	PermissionKickMembers:                      "Kick Members",
	PermissionBanMembers:                       "Ban Members",
	PermissionAdministrator:                    "Administrator",
	PermissionManageChannels:                   "Manage Channels",
	PermissionManageGuild:                      "Manage Server",
	PermissionAddReactions:                     "Add Reactions",
	PermissionViewAuditLog:                     "View Audit Logs",
	PermissionViewChannel:                      "View Channel",
	PermissionSendMessages:                     "Send Messages",
	PermissionSendTTSMessages:                  "Send TTS Messages",
	PermissionManageMessages:                   "Manage Messages",
	PermissionEmbedLinks:                       "Embed Links",
	PermissionAttachFiles:                      "Attach Files",
	PermissionReadMessageHistory:               "Read Message History",
	PermissionMentionEveryone:                  "Mention Everyone",
	PermissionUseExternalEmojis:                "Use External Emojis",
	PermissionConnect:                          "Connect",
	PermissionSpeak:                            "Speak",
	PermissionMuteMembers:                      "Mute Members",
	PermissionDeafenMembers:                    "Deafen Members",
	PermissionMoveMembers:                      "Move Members",
	PermissionUseVAD:                           "Use Voice Activity",
	PermissionPrioritySpeaker:                  "Priority Speaker",
	PermissionChangeNickname:                   "Change Nickname",
	PermissionManageNicknames:                  "Manage Nicknames",
	PermissionManageRoles:                      "Manage Roles",
	PermissionManageWebhooks:                   "Manage Webhooks",
	PermissionManageGuildExpressions:           "Manage Expressions",
	PermissionUseApplicationCommands:           "Use Application Commands",
	PermissionRequestToSpeak:                   "Request to Speak",
	PermissionManageEvents:                     "Manage Events",
	PermissionManageThreads:                    "Manage Threads",
	PermissionCreatePublicThreads:              "Create Public Threads",
	PermissionCreatePrivateThreads:             "Create Private Threads",
	PermissionUseExternalStickers:              "Use External Stickers",
	PermissionSendMessagesInThreads:            "Send Messages in Threads",
	PermissionUseEmbeddedActivities:            "Use Activities",
	PermissionModerateMembers:                  "Moderate Members",
	PermissionViewCreatorMonetizationAnalytics: "View Creator Monetization Analytics",
	PermissionUseSoundboard:                    "Use Soundboard",
	PermissionUseExternalSounds:                "Use External Sounds",
	PermissionStream:                           "Video",
	PermissionViewGuildInsights:                "View Server Insights",
	PermissionSendVoiceMessages:                "Send Voice Messages",
	PermissionSendPolls:                        "Create Polls",
	PermissionUseExternalApps:                  "Use External Apps",
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
	return flags.Add(p, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Permissions) Remove(bits ...Permissions) Permissions {
	return flags.Remove(p, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (p Permissions) Has(bits ...Permissions) bool {
	return flags.Has(p, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (p Permissions) Missing(bits ...Permissions) bool {
	return flags.Missing(p, bits...)
}
