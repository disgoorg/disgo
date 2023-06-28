package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var _ SoundboardSounds = (*soundboardSoundsImpl)(nil)

func NewSoundboardSounds(client Client) SoundboardSounds {
	return &soundboardSoundsImpl{client: client}
}

type SoundboardSounds interface {
	GetSoundboardDefaultSounds(opts ...RequestOpt) ([]discord.SoundboardSound, error)
	CreateGuildSoundboardSound(guildID snowflake.ID, soundCreate discord.SoundboardSoundCreate, opts ...RequestOpt) (*discord.SoundboardSound, error)
	UpdateGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, soundUpdate discord.SoundboardSoundUpdate, opts ...RequestOpt) (*discord.SoundboardSound, error)
	DeleteGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) error
}

type soundboardSoundsImpl struct {
	client Client
}

func (s *soundboardSoundsImpl) GetSoundboardDefaultSounds(opts ...RequestOpt) (sounds []discord.SoundboardSound, err error) {
	err = s.client.Do(GetSoundboardDefaultSounds.Compile(nil), nil, &sounds, opts...)
	return
}

func (s *soundboardSoundsImpl) CreateGuildSoundboardSound(guildID snowflake.ID, soundCreate discord.SoundboardSoundCreate, opts ...RequestOpt) (sound *discord.SoundboardSound, err error) {
	err = s.client.Do(CreateGuildSoundboardSound.Compile(nil, guildID), soundCreate, &sound, opts...)
	return
}

func (s *soundboardSoundsImpl) UpdateGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, soundUpdate discord.SoundboardSoundUpdate, opts ...RequestOpt) (sound *discord.SoundboardSound, err error) {
	err = s.client.Do(UpdateGuildSoundboardSound.Compile(nil, guildID, soundID), soundUpdate, &sound, opts...)
	return
}

func (s *soundboardSoundsImpl) DeleteGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) error {
	return s.client.Do(DeleteGuildSoundboardSound.Compile(nil, guildID, soundID), nil, nil, opts...)
}
