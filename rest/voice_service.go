package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ VoiceService = (*VoiceServiceImpl)(nil)

func NewVoiceService(restClient Client) VoiceService {
	return &VoiceServiceImpl{restClient: restClient}
}

type VoiceService interface {
	Service
	GetVoiceRegions(opts ...RequestOpt) ([]discord.VoiceRegion, Error)
}

type VoiceServiceImpl struct {
	restClient Client
}

func (s *VoiceServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *VoiceServiceImpl) GetVoiceRegions(opts ...RequestOpt) (regions []discord.VoiceRegion, rErr Error) {
	compiledRoute, err := route.GetVoiceRegions.Compile(nil)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, &regions, opts...)
	return
}
