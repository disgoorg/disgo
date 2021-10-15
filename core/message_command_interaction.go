package core

import "github.com/DisgoOrg/disgo/discord"

type MessageCommandInteractionFilter func(messageCommandInteraction *MessageCommandInteraction) bool

type MessageCommandInteraction struct {
	discord.MessageCommandInteraction
	CreateInteraction
	FollowupInteraction
	CommandID   discord.Snowflake
	CommandName string
	Resolved    *MessageCommandResolved
	TargetID    discord.Snowflake
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}

type MessageCommandResolved struct {
	Messages map[discord.Snowflake]*Message
}
