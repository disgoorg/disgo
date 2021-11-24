package discord

// GuildScheduledEvent a representation of a scheduled event in a Guild (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object)
type GuildScheduledEvent struct {
	ID                 Snowflake
	GuildID            Snowflake
	ChannelID          *Snowflake
	CreatorID          Snowflake
	Name               string
	Description        string
	ScheduledStartTime Time
	ScheduledEndTime   *Time
	PrivacyLevel       ScheduledEventPrivacyLevel
	Status             ScheduledEventStatus
	EntityType         ScheduledEventEntityType
	EntityID           *Snowflake
	EntityMetaData     *EntityMetaData
	Creator            User
	UserCount          int
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
