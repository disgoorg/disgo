package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewVoiceService(client Client) VoiceService {
	return nil
}

type VoiceService interface {
	Service
	GetVoiceRegions(opts ...rest.RequestOpt) []discord.VoiceRegion
}
