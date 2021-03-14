package internal

import (
	"github.com/DiscoOrg/disgo/api"
)

type EntityBuilder struct {
	Disgo api.Disgo
}

func (b EntityBuilder) createUser(user api.User) *api.User {
	user.Disgo = b.Disgo
	b.Disgo.Cache().CacheUser(&user)
	return &user
}

func (b EntityBuilder) createMessage(message api.Message) *api.Message {
	message.Disgo = b.Disgo
	//d.Disgo.Cache().CacheMessage(&message)
	return &message
}

func (b EntityBuilder) createMember(member api.Member) *api.Member {
	member.Disgo = b.Disgo
	b.Disgo.Cache().CacheMember(&member)
	return &member
}

func (b EntityBuilder) createDMChannel(channel api.DMChannel) *api.DMChannel {
	channel.Disgo = b.Disgo
	b.Disgo.Cache().CacheDMChannel(&channel)
	return &channel
}

func (b EntityBuilder) createTextChannel(channel api.TextChannel) *api.TextChannel {
	channel.MessageChannel.Disgo = b.Disgo
	b.Disgo.Cache().CacheTextChannel(&channel)
	return &channel
}

func (b EntityBuilder) createVoiceChannel(channel api.VoiceChannel) *api.VoiceChannel {
	channel.Disgo = b.Disgo
	b.Disgo.Cache().CacheVoiceChannel(&channel)
	return &channel
}

func (b EntityBuilder) createStoreChannel(channel api.StoreChannel) *api.StoreChannel {
	channel.Disgo = b.Disgo
	b.Disgo.Cache().CacheStoreChannel(&channel)
	return &channel
}
