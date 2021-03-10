package internal

import (
	"github.com/DiscoOrg/disgo/api"
)

type EntityBuilder struct {
	Disgo api.Disgo
}

func (d EntityBuilder) createMember() api.Member {
		return api.Member{}
}

func (d EntityBuilder) createDMChannel() api.DMChannel {
	return api.DMChannel{}
}

func (d EntityBuilder) createTextChannel() api.TextChannel {
	return api.TextChannel{}
}

func (d EntityBuilder) createVoiceChannel() api.VoiceChannel {
	return api.VoiceChannel{}
}

func (d EntityBuilder) createNewsChannel() api.NewsChannel {
	return api.NewsChannel{}
}

func (d EntityBuilder) createStoreChannel() api.StoreChannel {
	return api.StoreChannel{}
}
