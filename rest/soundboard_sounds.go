package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var _ SoundboardSounds = (*soundsImpl)(nil)

func NewSoundboardSounds(client Client) SoundboardSounds {
	return &soundsImpl{client: client}
}

type SoundboardSounds interface {
	GetSoundboardDefaultSounds(opts ...RequestOpt) ([]discord.SoundboardSound, error)
	GetGuildSoundboardSounds(guildID snowflake.ID, opts ...RequestOpt) ([]discord.SoundboardSound, error)
	CreateGuildSoundboardSound(guildID snowflake.ID, soundCreate discord.SoundboardSoundCreate, opts ...RequestOpt) (*discord.SoundboardSound, error)
	GetGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) (*discord.SoundboardSound, error)
	UpdateGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, soundUpdate discord.SoundboardSoundUpdate, opts ...RequestOpt) (*discord.SoundboardSound, error)
	DeleteGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) error
	SendSoundboardSound(channelID snowflake.ID, sendSound discord.SendSoundboardSound, opts ...RequestOpt) error
}

type soundsImpl struct {
	client Client
}

func (s *soundsImpl) GetSoundboardDefaultSounds(opts ...RequestOpt) (sounds []discord.SoundboardSound, err error) {
	err = s.client.Do(GetSoundboardDefaultSounds.Compile(nil), nil, &sounds, opts...)
	return
}

func (s *soundsImpl) GetGuildSoundboardSounds(guildID snowflake.ID, opts ...RequestOpt) (sounds []discord.SoundboardSound, err error) {
	var rs soundsResponse
	err = s.client.Do(GetGuildSoundboardSounds.Compile(nil, guildID), nil, &rs, opts...)
	if err == nil {
		sounds = rs.Items
	}
	return
}

func (s *soundsImpl) CreateGuildSoundboardSound(guildID snowflake.ID, soundCreate discord.SoundboardSoundCreate, opts ...RequestOpt) (sound *discord.SoundboardSound, err error) {
	err = s.client.Do(CreateGuildSoundboardSound.Compile(nil, guildID), soundCreate, &sound, opts...)
	return
}

func (s *soundsImpl) GetGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) (sound *discord.SoundboardSound, err error) {
	err = s.client.Do(GetGuildSoundboardSound.Compile(nil, guildID, soundID), nil, &sound, opts...)
	return
}

func (s *soundsImpl) UpdateGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, soundUpdate discord.SoundboardSoundUpdate, opts ...RequestOpt) (sound *discord.SoundboardSound, err error) {
	err = s.client.Do(UpdateGuildSoundboardSound.Compile(nil, guildID, soundID), soundUpdate, &sound, opts...)
	return
}

func (s *soundsImpl) DeleteGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuildSoundboardSound.Compile(nil, guildID, soundID), nil, nil, opts...)
}

func (s *soundsImpl) SendSoundboardSound(channelID snowflake.ID, sendSound discord.SendSoundboardSound, opts ...RequestOpt) error {
	return s.client.Do(SendSoundboardSound.Compile(nil, channelID), sendSound, nil, opts...)
}

type soundsResponse struct {
	Items []discord.SoundboardSound `json:"items"`
}
