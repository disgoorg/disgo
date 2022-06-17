package audio

import (
	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/audio/samplerate"
	"github.com/disgoorg/snowflake/v2"
)

func NewSampleRateReceiver(resampler *samplerate.Resampler, inputSampleRate int, outputSampleRate int, pcmFrameReceiver PCMFrameReceiver) PCMFrameReceiver {
	if resampler == nil {
		resampler = samplerate.CreateResampler(samplerate.ConverterTypeSincBestQuality, 2)
	}
	return &sampleRateReceiver{
		resampler:        resampler,
		pcmFrameReceiver: pcmFrameReceiver,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
	}
}

type sampleRateReceiver struct {
	resampler        *samplerate.Resampler
	pcmFrameReceiver PCMFrameReceiver
	inputSampleRate  int
	outputSampleRate int
}

func (p *sampleRateReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) error {
	newPCM := make([]int16, opus.GetOutputBuffSize(p.outputSampleRate, p.resampler.Channels()))
	var (
		inputFrames  int
		outputFrames int
	)
	if err := p.resampler.Process(packet.PCM, newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return err
	}

	packet.PCM = newPCM
	return p.pcmFrameReceiver.ReceivePCMFrame(userID, packet)
}

func (p *sampleRateReceiver) CleanupUser(userID snowflake.ID) {
	p.pcmFrameReceiver.CleanupUser(userID)
}

func (p *sampleRateReceiver) Close() {
	p.resampler.Destroy()
	p.pcmFrameReceiver.Close()
}
