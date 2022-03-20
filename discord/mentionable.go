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
	MentionTypeEmoji     = MentionType{regexp.MustCompile(`<a?:(\w+):(\d+)>`)}
	MentionTypeTimestamp = MentionType{regexp.MustCompile(`<t:(?P<time>-?\d{1,17})(?::(?P<format>[tTdDfFR]))?>`)}
	MentionTypeHere      = MentionType{regexp.MustCompile(`@here`)}
	MentionTypeEveryone  = MentionType{regexp.MustCompile(`@everyone`)}
)

type Mentionable interface {
	Mention() string
}

//goland:noinspection GoUnusedFunction
func ChannelMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

//goland:noinspection GoUnusedFunction
func UserTag(username string, discriminator string) string {
	return fmt.Sprintf("%s#%s", username, discriminator)
}

//goland:noinspection GoUnusedFunction
func UserMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@%s>", id)
}

//goland:noinspection GoUnusedFunction
func MemberMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@!%s>", id)
}

//goland:noinspection GoUnusedFunction
func RoleMention(id snowflake.Snowflake) string {
	return fmt.Sprintf("<@&%s>", id)
}

//goland:noinspection GoUnusedFunction
func EmojiMention(id snowflake.Snowflake, name string) string {
	return fmt.Sprintf("<:%s:%s>", name, id)
}

//goland:noinspection GoUnusedFunction
func AnimatedEmojiMention(id snowflake.Snowflake, name string) string {
	return fmt.Sprintf("<a:%s:%s>", name, id)
}

//goland:noinspection GoUnusedFunction
func TimestampMention(timestamp int64) string {
	return TimestampStyleNone.Format(timestamp)
}

//goland:noinspection GoUnusedFunction
func FormattedTimestampMention(timestamp int64, style TimestampStyle) string {
	return style.Format(timestamp)
}
