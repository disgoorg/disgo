package core

import "github.com/DisgoOrg/snowflake"

type (
	MessageFindFunc func(message *Message) bool

	MessageCache interface {
		Get(channelID snowflake.Snowflake, messageID snowflake.Snowflake) *Message
		GetCopy(channelID snowflake.Snowflake, messageID snowflake.Snowflake) *Message
		Set(message *Message) *Message
		Remove(channelID snowflake.Snowflake, messageID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Message
		All() map[snowflake.Snowflake][]*Message

		ChannelCache(channelID snowflake.Snowflake) map[snowflake.Snowflake]*Message
		ChannelAll(channelID snowflake.Snowflake) []*Message

		FindFirst(messageFindFunc MessageFindFunc) *Message
		FindAll(messageFindFunc MessageFindFunc) []*Message
	}

	messageCacheImpl struct {
		messageCachePolicy MessageCachePolicy
		messages           map[snowflake.Snowflake]map[snowflake.Snowflake]*Message
	}
)

func NewMessageCache(messageCachePolicy MessageCachePolicy) MessageCache {
	if messageCachePolicy == nil {
		messageCachePolicy = MessageCachePolicyDefault
	}
	return &messageCacheImpl{
		messageCachePolicy: messageCachePolicy,
		messages:           map[snowflake.Snowflake]map[snowflake.Snowflake]*Message{},
	}
}

func (c *messageCacheImpl) Get(channelID snowflake.Snowflake, messageID snowflake.Snowflake) *Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}
	return c.messages[channelID][messageID]
}

func (c *messageCacheImpl) GetCopy(channelID snowflake.Snowflake, messageID snowflake.Snowflake) *Message {
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
		c.messages[message.ChannelID] = map[snowflake.Snowflake]*Message{}
	}
	rol, ok := c.messages[message.ChannelID][message.ID]
	if ok {
		*rol = *message
		return rol
	}
	c.messages[message.ChannelID][message.ID] = message

	return message
}

func (c *messageCacheImpl) Remove(channelID snowflake.Snowflake, messageID snowflake.Snowflake) {
	if _, ok := c.messages[channelID]; !ok {
		return
	}
	delete(c.messages[channelID], messageID)
}

func (c *messageCacheImpl) Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Message {
	return c.messages
}

func (c *messageCacheImpl) All() map[snowflake.Snowflake][]*Message {
	messages := make(map[snowflake.Snowflake][]*Message, len(c.messages))
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

func (c *messageCacheImpl) ChannelCache(channelID snowflake.Snowflake) map[snowflake.Snowflake]*Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}
	return c.messages[channelID]
}

func (c *messageCacheImpl) ChannelAll(channelID snowflake.Snowflake) []*Message {
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
