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
	GetVoiceRegions(ctx context.Context) []discord.VoiceRegion
}
