package audio

import (
	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/audio/samplerate"
)

func NewSampleRateProvider(resampler *samplerate.Resampler, inputSampleRate int, outputSampleRate int, pcmFrameProvider PCMFrameProvider) PCMFrameProvider {
	if resampler == nil {
		resampler = samplerate.CreateResampler(samplerate.ConverterTypeSincBestQuality, 2)
	}
	return &sampleRateProvider{
		resampler:        resampler,
		pcmFrameProvider: pcmFrameProvider,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
	}
}

type sampleRateProvider struct {
	resampler        *samplerate.Resampler
	pcmFrameProvider PCMFrameProvider
	inputSampleRate  int
	outputSampleRate int
}

func (p *sampleRateProvider) ProvidePCMFrame() ([]int16, error) {
	pcm, err := p.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}

	newPCM := make([]int16, opus.GetOutputBuffSize(p.outputSampleRate, p.resampler.Channels()))
	var (
		inputFrames  int64
		outputFrames int64
	)
	if err = p.resampler.Process(pcm, newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return nil, err
	}

	return newPCM, nil
}

func (p *sampleRateProvider) Close() {
	p.resampler.Destroy()
	p.pcmFrameProvider.Close()
}
