package audio

import (
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

func (p *sampleRateCombinedReceiver) ReceiveCombinedPCMFrame(userIDs []snowflake.ID, packet *CombinedPCMPacket) {
	out := make([]int16, len(packet.PCM))
	var (
		inputFrames  int64
		outputFrames int64
	)
	if err := p.resampler.Process(packet.PCM, out, p.inputSampleRate, p.outputSampleRate, 1, &inputFrames, &outputFrames); err != nil {
		panic("ReceivePCMFrame: " + err.Error())
	}

	//fmt.Printf("sampleRateCombinedReceiver: inputFrames: %d, outputFrames: %d\n", inputFrames, outputFrames)

	packet.PCM = out[:outputFrames]
	p.pcmFrameReceiver.ReceiveCombinedPCMFrame(userIDs, packet)
}

func (p *sampleRateCombinedReceiver) Close() {
	p.resampler.Destroy()
	p.pcmFrameReceiver.Close()
}
