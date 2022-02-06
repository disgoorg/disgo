package discord

import (
	"fmt"
	"regexp"

	"github.com/DisgoOrg/snowflake"
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
func channelMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

//goland:noinspection GoUnusedFunction
func userTag(username string, discriminator string) string {
	return fmt.Sprintf("%s#%s", username, discriminator)
}

//goland:noinspection GoUnusedFunction
func userMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@%s>", id)
}

//goland:noinspection GoUnusedFunction
func memberMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@!%s>", id)
}

//goland:noinspection GoUnusedFunction
func roleMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@&%s>", id)
}

//goland:noinspection GoUnusedFunction
func emojiMention(id snowflake.Snowflake, name string) string {
	return fmt.Sprintf("<:%s:%s>", name, id)
}

//goland:noinspection GoUnusedFunction
func animatedEmojiMention(id snowflake.Snowflake, name string) string {
	return fmt.Sprintf("<a:%s:%s>", name, id)
}

//goland:noinspection GoUnusedFunction
func timestampMention(timestamp int64) string {
	return fmt.Sprintf("<t:%d>", timestamp)
}

//goland:noinspection GoUnusedFunction
func formattedTimestampMention(timestamp int64, style TimestampStyle) string {
	return fmt.Sprintf("<t:%d:%s>", timestamp, style)
}
