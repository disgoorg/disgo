package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
)

var _ Voice = (*voiceImpl)(nil)

func NewVoice(restClient Client) Voice {
	return &voiceImpl{restClient: restClient}
}

type Voice interface {
	GetVoiceRegions(opts ...RequestOpt) ([]discord.VoiceRegion, error)
}

type voiceImpl struct {
	restClient Client
}

func (s *voiceImpl) GetVoiceRegions(opts ...RequestOpt) (regions []discord.VoiceRegion, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetVoiceRegions.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &regions, opts...)
	return
}
