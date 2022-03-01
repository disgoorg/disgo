package core

import "github.com/DisgoOrg/snowflake"

type MemberChunkingFilter func(guildID snowflake.Snowflake) bool

//goland:noinspection GoUnusedGlobalVariable
var (
	MemberChunkingFilterAll  MemberChunkingFilter = func(_ snowflake.Snowflake) bool { return true }
	MemberChunkingFilterNone MemberChunkingFilter = func(_ snowflake.Snowflake) bool { return false }
)

func (f MemberChunkingFilter) Include(guildIDs ...snowflake.Snowflake) MemberChunkingFilter {
	return func(guildID snowflake.Snowflake) bool {
		return f(guildID) || func(guildID snowflake.Snowflake) bool {
			for _, gID := range guildIDs {
				if gID == guildID {
					return true
				}
			}
			return false
		}(guildID)
	}
}

func (f MemberChunkingFilter) Exclude(guildIDs ...snowflake.Snowflake) MemberChunkingFilter {
	return func(guildID snowflake.Snowflake) bool {
		return f(guildID) || func(guildID snowflake.Snowflake) bool {
			for _, gID := range guildIDs {
				if gID == guildID {
					return false
				}
			}
			return true
		}(guildID)
	}
}
