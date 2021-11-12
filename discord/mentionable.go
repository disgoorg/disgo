package discord

import "fmt"

type Mentionable interface {
	Mention() string
}

func channelMention(id Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

func userMention(id Snowflake) string {
	return fmt.Sprintf("<@%s>", id)
}

func memberMention(id Snowflake) string {
	return fmt.Sprintf("<@!%s>", id)
}

func roleMention(id Snowflake) string {
	return fmt.Sprintf("<@&%s>", id)
}

func emojiMention(id Snowflake, name string) string {
	return fmt.Sprintf("<:%s:%s>", id, name)
}

func animatedEmojiMention(id Snowflake, name string) string {
	return fmt.Sprintf("<a:%s:%s>", id, name)
}

func timestampMention(timestamp int64) string {
	return fmt.Sprintf("<t:%d>", timestamp)
}

func formattedTimestampMention(timestamp int64, style TimestampStyle) string {
	return fmt.Sprintf("<t:%d:%s>", timestamp, style)
}
