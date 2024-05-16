package bot

import (
	"slices"

	"github.com/disgoorg/snowflake/v2"
)

// MemberChunkingFilterAll is a MemberChunkingFilter which includes all guilds.
func MemberChunkingFilterAll(_ snowflake.ID) bool { return true }

// MemberChunkingFilterNone is a MemberChunkingFilter which excludes all guilds.
func MemberChunkingFilterNone(_ snowflake.ID) bool { return false }

// MemberChunkingFilterDefault is the default MemberChunkingFilter.
func MemberChunkingFilterDefault(guildID snowflake.ID) bool {
	return MemberChunkingFilterNone(guildID)
}

// MemberChunkingFilterIncludeGuildIDs returns a MemberChunkingFilter which includes the given guildIDs.
func MemberChunkingFilterIncludeGuildIDs(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return slices.Contains(guildIDs, guildID)
	}
}

// MemberChunkingFilterExcludeGuildIDs returns a MemberChunkingFilter which excludes the given guildIDs.
func MemberChunkingFilterExcludeGuildIDs(guildIDs ...snowflake.ID) MemberChunkingFilter {
	return func(guildID snowflake.ID) bool {
		return slices.Contains(guildIDs, guildID)
	}
}

// MemberChunkingFilter is a filter that can be used to filter from which guilds to request members from.
type MemberChunkingFilter func(guildID snowflake.ID) bool

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

// AnyMemberChunkingFilter is a shorthand for MemberChunkingFilter.Or(MemberChunkingFilter).Or(MemberChunkingFilter) etc.
func AnyMemberChunkingFilter(filters ...MemberChunkingFilter) MemberChunkingFilter {
	var filter MemberChunkingFilter
	for _, f := range filters {
		if filter == nil {
			filter = f
			continue
		}
		filter = filter.Or(f)
	}
	return filter
}

// AllMemberChunkingFilters is a shorthand for MemberChunkingFilter.And(MemberChunkingFilter).And(MemberChunkingFilter) etc.
func AllMemberChunkingFilters(filters ...MemberChunkingFilter) MemberChunkingFilter {
	var filter MemberChunkingFilter
	for _, f := range filters {
		if filter == nil {
			filter = f
			continue
		}
		filter = filter.And(f)
	}
	return filter
}
