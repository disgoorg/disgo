package discord

import (
	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

type ThreadCreateFromMessage struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    int                 `json:"rate_limit_per_user,omitempty"`
}

type ThreadChannelPostCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    int                 `json:"rate_limit_per_user,omitempty"`
	Message             MessageCreate       `json:"message"`
	AppliedTags         []snowflake.ID      `json:"applied_tags,omitempty"`
}

func (c ThreadChannelPostCreate) ToBody() (any, error) {
	if len(c.Message.Files) > 0 {
		c.Message.Attachments = parseAttachments(c.Message.Files)
		return PayloadWithFiles(c, c.Message.Files...)
	}
	return c, nil
}

type ThreadChannelPost struct {
	GuildThread
	Message Message `json:"message"`
}

func (c *ThreadChannelPost) UnmarshalJSON(data []byte) error {
	var thread GuildThread
	if err := json.Unmarshal(data, &thread); err != nil {
		return err
	}

	c.GuildThread = thread

	var v struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.GuildThread = thread
	c.Message = v.Message
	return nil
}

func (c ThreadChannelPost) MarshalJSON() ([]byte, error) {
	data1, err := json.Marshal(c.GuildThread)
	if err != nil {
		return nil, err
	}

	data2, err := json.Marshal(struct {
		Message Message `json:"message"`
	}{
		Message: c.Message,
	})
	if err != nil {
		return nil, err
	}

	return json.Merge(data1, data2)
}

type ThreadCreate interface {
	json.Marshaler
	Type() ChannelType
}

type GuildNewsThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
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

func (GuildNewsThreadCreate) Type() ChannelType {
	return ChannelTypeGuildNewsThread
}

type GuildPublicThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
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

func (GuildPublicThreadCreate) Type() ChannelType {
	return ChannelTypeGuildPublicThread
}

type GuildPrivateThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
	Invitable           *bool               `json:"invitable,omitempty"`
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

func (GuildPrivateThreadCreate) Type() ChannelType {
	return ChannelTypeGuildPrivateThread
}

type GetThreads struct {
	Threads []GuildThread  `json:"threads"`
	Members []ThreadMember `json:"members"`
	HasMore bool           `json:"has_more"`
}

type GetAllThreads struct {
	Threads []GuildThread  `json:"threads"`
	Members []ThreadMember `json:"members"`
}
