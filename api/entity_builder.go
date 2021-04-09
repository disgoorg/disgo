package api

type EntityBuilder interface {
	Disgo() Disgo

	CreateTextChannel(channel *Channel) *TextChannel
	CreateVoiceChannel(channel *Channel) *VoiceChannel
	CreateStoreChannel(channel *Channel) *StoreChannel
	CreateCategory(channel *Channel) *Category
	CreateDMChannel(channel *Channel) *DMChannel
}
