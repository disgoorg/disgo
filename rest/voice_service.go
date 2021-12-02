package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service      = (*voiceServiceImpl)(nil)
	_ VoiceService = (*voiceServiceImpl)(nil)
)

func NewVoiceService(restClient Client) VoiceService {
	return &voiceServiceImpl{restClient: restClient}
}

type VoiceService interface {
	Service
	GetVoiceRegions(opts ...RequestOpt) ([]discord.VoiceRegion, error)
}

type voiceServiceImpl struct {
	restClient Client
}

func (s *voiceServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *voiceServiceImpl) GetVoiceRegions(opts ...RequestOpt) (regions []discord.VoiceRegion, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetVoiceRegions.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &regions, opts...)
	return
}
