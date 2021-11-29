package discord

// GuildScheduledEvent a representation of a scheduled event in a Guild (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object)
type GuildScheduledEvent struct {
	ID                 Snowflake                  `json:"id"`
	GuildID            Snowflake                  `json:"guild_id"`
	ChannelID          *Snowflake                 `json:"channel_id"`
	CreatorID          Snowflake                  `json:"creator_id"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	ScheduledStartTime Time                       `json:"scheduled_start_time"`
	ScheduledEndTime   *Time                      `json:"scheduled_end_time"`
	PrivacyLevel       ScheduledEventPrivacyLevel `json:"privacy_level"`
	Status             ScheduledEventStatus       `json:"status"`
	EntityType         ScheduledEventEntityType   `json:"entity_type"`
	EntityID           *Snowflake                 `json:"entity_id"`
	EntityMetaData     *EntityMetaData            `json:"entity_metadata"`
	Creator            User                       `json:"creator"`
	UserCount          int                        `json:"user_count"`
}

type GuildScheduledEventCreate struct {
	ChannelID          Snowflake                  `json:"channel_id,omitempty"`
	EntityMetaData     *EntityMetaData            `json:"entity_metadata,omitempty"`
	Name               string                     `json:"name"`
	PrivacyLevel       ScheduledEventPrivacyLevel `json:"privacy_level"`
	ScheduledStartTime Time                       `json:"scheduled_start_time"`
	ScheduledEndTime   Time                       `json:"scheduled_end_time,omitempty"`
	Description        string                     `json:"description,omitempty"`
	EntityType         ScheduledEventEntityType   `json:"entity_type"`
}

type GuildScheduledEventUpdate struct {
	ChannelID          *Snowflake                  `json:"channel_id,omitempty"`
	EntityMetaData     *EntityMetaData             `json:"entity_metadata,omitempty"`
	Name               string                      `json:"name,omitempty"`
	PrivacyLevel       *ScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	ScheduledStartTime *Time                       `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *Time                       `json:"scheduled_end_time,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	EntityType         *ScheduledEventEntityType   `json:"entity_type,omitempty"`
}

type GuildScheduledEventUser struct {
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id"`
	User                  User      `json:"user"`
	Member                *Member   `json:"member"`
}

// ScheduledEventPrivacyLevel the privacy level of the ScheduledEventPrivacyLevel (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level)
type ScheduledEventPrivacyLevel int

//goland:noinspection GoUnusedConst
const (
	_ ScheduledEventPrivacyLevel = iota + 1
	ScheduledEventPrivacyLevelGuildOnly
)

// ScheduledEventStatus the status of the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status)
type ScheduledEventStatus int

//goland:noinspection GoUnusedConst
const (
	ScheduledEventStatusScheduled ScheduledEventStatus = iota + 1
	ScheduledEventStatusActive
	ScheduledEventStatusCompleted
	ScheduledEventStatusCancelled
)

// ScheduledEventEntityType the type of the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types)
type ScheduledEventEntityType int

//goland:noinspection GoUnusedConst
const (
	ScheduledEventEntityTypeStageInstance ScheduledEventEntityType = iota + 1
	ScheduledEventEntityTypeVoice
	ScheduledEventEntityTypeExternal
)

// EntityMetaData additional metadata for the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata)
type EntityMetaData struct {
	Location string
}
