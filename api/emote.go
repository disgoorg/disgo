package api

// An Emote allows you to interact with custom emojis in discord.
type Emote struct {
	ID       Snowflake
	Name     string
	Animated bool
}

// Mention returns the string used to send the emoji
func (e Emote) Mention() string {
	start := "<:"
	if e.Animated {
		start = "<a:"
	}
	return start + e.Name + ":" + e.ID.String() + ">"
}

func (e Emote) String() string {
	return e.Mention()
}

// Reaction returns the identifier used for adding and removing reactions for messages in discord
func (e Emote) Reaction() string {
	return ":" + e.Name + ":" + e.ID.String()
}
