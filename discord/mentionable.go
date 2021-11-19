package discord

import "fmt"

type Mentionable interface {
	Mention() string
}

//goland:noinspection GoUnusedFunction
func channelMention(id Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
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
