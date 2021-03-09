package internal

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/models"
)

type EntityBuilder struct {
	Disgo api.Disgo
}

func (d EntityBuilder) createMember() models.Member {
		return models.Member{}
}

func (d EntityBuilder) createDMChannel() models.DMChannel {
	return models.DMChannel{}
}

func (d EntityBuilder) createTextChannel() models.TextChannel {
	return models.TextChannel{}
}

func (d EntityBuilder) createVoiceChannel() models.VoiceChannel {
	return models.VoiceChannel{}
}

func (d EntityBuilder) createNewsChannel() models.NewsChannel {
	return models.NewsChannel{}
}

func (d EntityBuilder) createStoreChannel() models.StoreChannel {
	return models.StoreChannel{}
}
