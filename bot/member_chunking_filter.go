package bot

import (
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/exp/slices"
)

// MemberChunkingFilter is a filter that can be used to filter from which guilds to request members from.
type MemberChunkingFilter func(guildID snowflake.ID) bool

var (
	MemberChunkingFilterAll  MemberChunkingFilter = func(_ snowflake.ID) bool { return true }
	MemberChunkingFilterNone MemberChunkingFilter = func(_ snowflake.ID) bool { return false }
)

func MemberChunkingFilterIncludeGuildIDs(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return slices.Contains(guildIDs, guildID)
	}
}

func MemberChunkingFilterExcludeGuildIDs(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return slices.Contains(guildIDs, guildID)
	}
}

// Or allows you to combine the MemberChunkingFilter with another, meaning either of them needs to be true for the guild to be chunked.
func (f MemberChunkingFilter) Or(filter MemberChunkingFilter) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return f(guildID) || filter(guildID)
	}
}

// And allows you to require both MemberChunkingFilter(s) to be true for the guild to be chunked
func (f MemberChunkingFilter) And(filter MemberChunkingFilter) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return f(guildID) && filter(guildID)
	}
}
