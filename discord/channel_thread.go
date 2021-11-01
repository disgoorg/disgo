package discord

type ChannelGuildThread interface {
	guildThread()
}

type GuildNewsThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

type GuildPublicThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

type GuildPrivateThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

type ThreadMetadata struct {
	Archived            bool                `json:"archived"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	ArchiveTimestamp    Time                `json:"archive_timestamp"`
	Locked              bool                `json:"locked"`
	Invitable           bool                `json:"invitable"`
}
