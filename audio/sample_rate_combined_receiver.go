package audio

import (
	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/audio/samplerate"
	"github.com/disgoorg/snowflake/v2"
)

func NewSampleRateCombinedReceiver(resampler *samplerate.Resampler, inputSampleRate int, outputSampleRate int, pcmFrameReceiver PCMCombinedFrameReceiver) PCMCombinedFrameReceiver {
	if resampler == nil {
		resampler = samplerate.CreateResampler(samplerate.ConverterTypeSincBestQuality, 2)
	}
	return &sampleRateCombinedReceiver{
		resampler:        resampler,
		pcmFrameReceiver: pcmFrameReceiver,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
	}
}

type sampleRateCombinedReceiver struct {
	resampler        *samplerate.Resampler
	pcmFrameReceiver PCMCombinedFrameReceiver
	inputSampleRate  int
	outputSampleRate int
}

func (p *sampleRateCombinedReceiver) ReceiveCombinedPCMFrame(userIDs []snowflake.ID, packet *CombinedPCMPacket) error {
	newPCM := make([]int16, opus.GetOutputBuffSize(p.outputSampleRate, p.resampler.Channels()))
	var (
		inputFrames  int
		outputFrames int
	)
	if err := p.resampler.Process(packet.PCM, newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return err
	}

	packet.PCM = newPCM
	return p.pcmFrameReceiver.ReceiveCombinedPCMFrame(userIDs, packet)
}

func (p *sampleRateCombinedReceiver) Close() {
	p.resampler.Destroy()
	p.pcmFrameReceiver.Close()
}
