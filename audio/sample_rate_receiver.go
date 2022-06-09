package audio

import (
	"fmt"

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

func (p *sampleRateReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) {
	out := make([]int16, len(packet.PCM))
	var (
		inputFrames  int64
		outputFrames int64
	)
	if err := p.resampler.Process(packet.PCM, out, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		panic("ReceivePCMFrame: " + err.Error())
	}

	fmt.Printf("sampleRateReceiver: inputFrames: %d, outputFrames: %d\n", inputFrames, outputFrames)

	packet.PCM = out
	p.pcmFrameReceiver.ReceivePCMFrame(userID, packet)
}

func (p *sampleRateReceiver) CleanupUser(userID snowflake.ID) {
	p.pcmFrameReceiver.CleanupUser(userID)
}

func (p *sampleRateReceiver) Close() {
	p.resampler.Destroy()
	p.pcmFrameReceiver.Close()
}
