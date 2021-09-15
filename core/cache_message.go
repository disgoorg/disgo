package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	MessageFindFunc func(message *Message) bool

	MessageCache interface {
		Get(channelID discord.Snowflake, messageID discord.Snowflake) *Message
		GetCopy(channelID discord.Snowflake, messageID discord.Snowflake) *Message
		Set(message *Message) *Message
		Remove(channelID discord.Snowflake, messageID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Message
		All() map[discord.Snowflake][]*Message

		ChannelCache(channelID discord.Snowflake) map[discord.Snowflake]*Message
		ChannelAll(channelID discord.Snowflake) []*Message

		FindFirst(messageFindFunc MessageFindFunc) *Message
		FindAll(messageFindFunc MessageFindFunc) []*Message
	}

	messageCacheImpl struct {
		messageCachePolicy MessageCachePolicy
		messages           map[discord.Snowflake]map[discord.Snowflake]*Message
	}
)

func NewMessageCache(messageCachePolicy MessageCachePolicy) MessageCache {
	if messageCachePolicy == nil {
		messageCachePolicy = MessageCachePolicyDefault
	}
	return &messageCacheImpl{
		messageCachePolicy: messageCachePolicy,
		messages:           map[discord.Snowflake]map[discord.Snowflake]*Message{},
	}
}

func (c *messageCacheImpl) Get(channelID discord.Snowflake, messageID discord.Snowflake) *Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}
	return c.messages[channelID][messageID]
}

func (c *messageCacheImpl) GetCopy(channelID discord.Snowflake, messageID discord.Snowflake) *Message {
	if message := c.Get(channelID, messageID); message != nil {
		me := *message
		return &me
	}
	return nil
}

func (c *messageCacheImpl) Set(message *Message) *Message {
	if !c.messageCachePolicy(message) {
		return message
	}
	if _, ok := c.messages[message.ChannelID]; !ok {
		c.messages[message.ChannelID] = map[discord.Snowflake]*Message{}
	}
	rol, ok := c.messages[message.ChannelID][message.ID]
	if ok {
		*rol = *message
		return rol
	}
	c.messages[message.ChannelID][message.ID] = message

	return message
}

func (c *messageCacheImpl) Remove(channelID discord.Snowflake, messageID discord.Snowflake) {
	if _, ok := c.messages[channelID]; !ok {
		return
	}
	delete(c.messages[channelID], messageID)
}

func (c *messageCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Message {
	return c.messages
}

func (c *messageCacheImpl) All() map[discord.Snowflake][]*Message {
	messages := make(map[discord.Snowflake][]*Message, len(c.messages))
	for channelID, channelMessages := range c.messages {
		messages[channelID] = make([]*Message, len(channelMessages))
		i := 0
		for _, channelMessage := range channelMessages {
			messages[channelID] = append(messages[channelID], channelMessage)
		}
		i++
	}
	return messages
}

func (c *messageCacheImpl) ChannelCache(channelID discord.Snowflake) map[discord.Snowflake]*Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}
	return c.messages[channelID]
}

func (c *messageCacheImpl) ChannelAll(channelID discord.Snowflake) []*Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}
	messages := make([]*Message, len(c.messages[channelID]))
	i := 0
	for _, message := range c.messages[channelID] {
		messages = append(messages, message)
		i++
	}
	return messages
}

func (c *messageCacheImpl) FindFirst(messageFindFunc MessageFindFunc) *Message {
	for _, channelMessages := range c.messages {
		for _, message := range channelMessages {
			if messageFindFunc(message) {
				return message
			}
		}
	}
	return nil
}

func (c *messageCacheImpl) FindAll(messageFindFunc MessageFindFunc) []*Message {
	var messages []*Message
	for _, channelMessages := range c.messages {
		for _, message := range channelMessages {
			if messageFindFunc(message) {
				messages = append(messages, message)
			}
		}
	}
	return messages
}
