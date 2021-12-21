package core

import "github.com/DisgoOrg/disgo/discord"

type MemberChunkingFilter func(guildID discord.Snowflake) bool

//goland:noinspection GoUnusedGlobalVariable
var (
	MemberChunkingFilterAll  MemberChunkingFilter = func(_ discord.Snowflake) bool { return true }
	MemberChunkingFilterNone MemberChunkingFilter = func(_ discord.Snowflake) bool { return false }
)

func (f MemberChunkingFilter) Include(guildIDs ...discord.Snowflake) MemberChunkingFilter {
	return func(guildID discord.Snowflake) bool {
		return f(guildID) || func(guildID discord.Snowflake) bool {
			for _, gID := range guildIDs {
				if gID == guildID {
					return true
				}
			}
			return false
		}(guildID)
	}
}

func (f MemberChunkingFilter) Exclude(guildIDs ...discord.Snowflake) MemberChunkingFilter {
	return func(guildID discord.Snowflake) bool {
		return f(guildID) || func(guildID discord.Snowflake) bool {
			for _, gID := range guildIDs {
				if gID == guildID {
					return false
				}
			}
			return true
		}(guildID)
	}
}
