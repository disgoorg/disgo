package audio

import (
	"fmt"

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

func (p *sampleRateProvider) ProvidePCMFrame() []int16 {
	pcm := p.pcmFrameProvider.ProvidePCMFrame()

	newPCM := make([]int16, len(pcm))
	var (
		inputFrames  int64
		outputFrames int64
	)
	if err := p.resampler.Process(pcm, newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		panic("sampleRateProvider: ReceivePCMFrame: " + err.Error())
	}

	fmt.Printf("inputFrames: %d, outputFrames: %d\n", inputFrames, outputFrames)

	return newPCM
}

func (p *sampleRateProvider) Close() {
	p.resampler.Destroy()
	p.pcmFrameProvider.Close()
}
