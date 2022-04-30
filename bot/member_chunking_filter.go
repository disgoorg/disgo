package bot

import "github.com/disgoorg/snowflake"

// MemberChunkingFilter is a filter that can be used to filter from which guilds to request members from.
type MemberChunkingFilter func(guildID snowflake.Snowflake) bool

var (
	MemberChunkingFilterAll  MemberChunkingFilter = func(_ snowflake.Snowflake) bool { return true }
	MemberChunkingFilterNone MemberChunkingFilter = func(_ snowflake.Snowflake) bool { return false }
)

// Include includes the given guilds from being chunked.
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

// Exclude excludes the given guilds from being chunked.
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
