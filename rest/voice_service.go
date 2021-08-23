package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewVoiceService(client Client) VoiceService {
	return nil
}

type VoiceService interface {
	Service
	GetVoiceRegions(opts ...RequestOpt) []discord.VoiceRegion
}
