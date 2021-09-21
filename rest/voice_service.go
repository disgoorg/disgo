package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ VoiceService = (*voiceServiceImpl)(nil)

func NewVoiceService(restClient Client) VoiceService {
	return &voiceServiceImpl{restClient: restClient}
}

type VoiceService interface {
	Service
	GetVoiceRegions(opts ...RequestOpt) ([]discord.VoiceRegion, Error)
}

type voiceServiceImpl struct {
	restClient Client
}

func (s *voiceServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *voiceServiceImpl) GetVoiceRegions(opts ...RequestOpt) (regions []discord.VoiceRegion, rErr Error) {
	compiledRoute, err := route.GetVoiceRegions.Compile(nil)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, &regions, opts...)
	return
}
