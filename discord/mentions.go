package discord

import "regexp"

type MentionType = *regexp.Regexp

var (
	MentionTypeUser          = regexp.MustCompile("<@!?(\\d+)>")
	MentionTypeRole          = regexp.MustCompile("<@&(\\d+)>")
	MentionTypeChannel       = regexp.MustCompile("<#(\\d+)>")
	MentionTypeEmoji         = regexp.MustCompile("<a?:([a-zA-Z0-9_]+):([0-9]+)>")
	MentionTypeTimestamp     = regexp.MustCompile("<t:([0-9]+)(:[tTdDfFR])?>")
	MentionTypeHere          = regexp.MustCompile("@here")
	MentionTypeEveryone      = regexp.MustCompile("@everyone")
)
