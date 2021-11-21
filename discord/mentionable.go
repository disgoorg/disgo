package discord

import (
	"fmt"
	"regexp"
)

type MentionType struct {
	*regexp.Regexp
}

//goland:noinspection GoUnusedGlobalVariable
var (
	MentionTypeUser      = MentionType{regexp.MustCompile(`<@!?(\d+)>`)}
	MentionTypeRole      = MentionType{regexp.MustCompile(`<@&(\d+)>`)}
	MentionTypeChannel   = MentionType{regexp.MustCompile(`<#(\d+)>`)}
	MentionTypeEmoji     = MentionType{regexp.MustCompile(`<a?:([a-zA-Z0-9_]+):([0-9]+)>`)}
	MentionTypeTimestamp = MentionType{regexp.MustCompile(`<t:(?P<time>-?\d{1,17})(?::(?P<format>[tTdDfFR]))?>`)}
	MentionTypeHere      = MentionType{regexp.MustCompile(`@here`)}
	MentionTypeEveryone  = MentionType{regexp.MustCompile(`@everyone`)}
)

type Mentionable interface {
	Mention() string
}

//goland:noinspection GoUnusedFunction
func channelMention(id Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

//goland:noinspection GoUnusedFunction
func userTag(username string, discriminator string) string {
	return fmt.Sprintf("%s#%s", username, discriminator)
}

//goland:noinspection GoUnusedFunction
func userMention(id Snowflake) string {
	return fmt.Sprintf("<@%s>", id)
}

//goland:noinspection GoUnusedFunction
func memberMention(id Snowflake) string {
	return fmt.Sprintf("<@!%s>", id)
}

//goland:noinspection GoUnusedFunction
func roleMention(id Snowflake) string {
	return fmt.Sprintf("<@&%s>", id)
}

//goland:noinspection GoUnusedFunction
func emojiMention(id Snowflake, name string) string {
	return fmt.Sprintf("<:%s:%s>", id, name)
}

//goland:noinspection GoUnusedFunction
func animatedEmojiMention(id Snowflake, name string) string {
	return fmt.Sprintf("<a:%s:%s>", id, name)
}

//goland:noinspection GoUnusedFunction
func timestampMention(timestamp int64) string {
	return fmt.Sprintf("<t:%d>", timestamp)
}

//goland:noinspection GoUnusedFunction
func formattedTimestampMention(timestamp int64, style TimestampStyle) string {
	return fmt.Sprintf("<t:%d:%s>", timestamp, style)
}
