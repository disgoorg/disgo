package core

import "github.com/DisgoOrg/disgo/discord"

type MessageCommandInteractionFilter func(messageCommandInteraction *MessageCommandInteraction) bool

type MessageCommandInteraction struct {
	*InteractionFields
	CommandID   discord.Snowflake
	CommandName string
	Resolved    *MessageCommandResolved
	TargetID    discord.Snowflake
}

func (i *MessageCommandInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeApplicationCommand
}

func (i *MessageCommandInteraction) ApplicationCommandType() discord.ApplicationCommandType {
	return discord.ApplicationCommandTypeMessage
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}

type MessageCommandResolved struct {
	Messages map[discord.Snowflake]*Message
}
