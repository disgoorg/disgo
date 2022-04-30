package bot

import "github.com/disgoorg/snowflake/v2"

// MemberChunkingFilter is a filter that can be used to filter from which guilds to request members from.
type MemberChunkingFilter func(guildID snowflake.ID) bool

var (
	MemberChunkingFilterAll  MemberChunkingFilter = func(_ snowflake.ID) bool { return true }
	MemberChunkingFilterNone MemberChunkingFilter = func(_ snowflake.ID) bool { return false }
)

// Include includes the given guilds from being chunked.
func (f MemberChunkingFilter) Include(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return f(guildID) || func(guildID snowflake.ID) bool {
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
func (f MemberChunkingFilter) Exclude(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return f(guildID) || func(guildID snowflake.ID) bool {
			for _, gID := range guildIDs {
				if gID == guildID {
					return false
				}
			}
			return true
		}(guildID)
	}
}
