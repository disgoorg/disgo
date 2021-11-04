package discord

import "github.com/DisgoOrg/disgo/json"

type ThreadCreateWithMessage struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
}

type ThreadCreate interface {
	json.Marshaler
	Type() ChannelType
}

type GuildNewsThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	Invitable           bool                `json:"invitable"`
}

func (c GuildNewsThreadCreate) MarshalJSON() ([]byte, error) {
	type guildNewsThreadCreate GuildNewsThreadCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildNewsThreadCreate
	}{
		Type:                  c.Type(),
		guildNewsThreadCreate: guildNewsThreadCreate(c),
	})
}

func (_ GuildNewsThreadCreate) Type() ChannelType {
	return ChannelTypeGuildNewsThread
}

type GuildPublicThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	Invitable           bool                `json:"invitable"`
}

func (c GuildPublicThreadCreate) MarshalJSON() ([]byte, error) {
	type guildPublicThreadCreate GuildPublicThreadCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildPublicThreadCreate
	}{
		Type:                    c.Type(),
		guildPublicThreadCreate: guildPublicThreadCreate(c),
	})
}

func (_ GuildPublicThreadCreate) Type() ChannelType {
	return ChannelTypeGuildPublicThread
}

type GuildPrivateThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	Invitable           bool                `json:"invitable"`
}

func (c GuildPrivateThreadCreate) MarshalJSON() ([]byte, error) {
	type guildPrivateThreadCreate GuildPrivateThreadCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildPrivateThreadCreate
	}{
		Type:                     c.Type(),
		guildPrivateThreadCreate: guildPrivateThreadCreate(c),
	})
}

func (_ GuildPrivateThreadCreate) Type() ChannelType {
	return ChannelTypeGuildPrivateThread
}

type GetThreads struct {
	Threads []GuildThread `json:"threads"`
	Members []Member      `json:"members"`
	HasMore bool          `json:"has_more"`
}
