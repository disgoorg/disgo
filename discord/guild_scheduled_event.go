package discord

import (
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

// GuildScheduledEvent a representation of a scheduled event in a Guild (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object)
type GuildScheduledEvent struct {
	ID                 snowflake.ID                  `json:"id"`
	GuildID            snowflake.ID                  `json:"guild_id"`
	ChannelID          *snowflake.ID                 `json:"channel_id"`
	CreatorID          snowflake.ID                  `json:"creator_id"`
	Name               string                        `json:"name"`
	Description        string                        `json:"description"`
	ScheduledStartTime time.Time                     `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                    `json:"scheduled_end_time"`
	PrivacyLevel       ScheduledEventPrivacyLevel    `json:"privacy_level"`
	Status             ScheduledEventStatus          `json:"status"`
	EntityType         ScheduledEventEntityType      `json:"entity_type"`
	EntityID           *snowflake.ID                 `json:"entity_id"`
	EntityMetaData     *EntityMetaData               `json:"entity_metadata"`
	Creator            User                          `json:"creator"`
	UserCount          int                           `json:"user_count"`
	Image              *string                       `json:"image"`
	RecurrenceRule     *ScheduledEventRecurrenceRule `json:"recurrence_rule"`
}

func (e GuildScheduledEvent) CreatedAt() time.Time {
	return e.ID.Time()
}

// CoverURL returns the cover URL if set or nil
func (e GuildScheduledEvent) CoverURL(opts ...CDNOpt) *string {
	if e.Image == nil {
		return nil
	}
	url := formatAssetURL(GuildScheduledEventCover, opts, e.ID, e.Image)
	return &url
}

type GuildScheduledEventCreate struct {
	ChannelID          snowflake.ID                  `json:"channel_id,omitempty"`
	EntityMetaData     *EntityMetaData               `json:"entity_metadata,omitempty"`
	Name               string                        `json:"name"`
	PrivacyLevel       ScheduledEventPrivacyLevel    `json:"privacy_level"`
	ScheduledStartTime time.Time                     `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                    `json:"scheduled_end_time,omitempty"`
	Description        string                        `json:"description,omitempty"`
	EntityType         ScheduledEventEntityType      `json:"entity_type"`
	Image              *Icon                         `json:"image,omitempty"`
	RecurrenceRule     *ScheduledEventRecurrenceRule `json:"recurrence_rule,omitempty"`
}

type GuildScheduledEventUpdate struct {
	ChannelID          *snowflake.ID                                `json:"channel_id,omitempty"`
	EntityMetaData     *EntityMetaData                              `json:"entity_metadata,omitempty"`
	Name               string                                       `json:"name,omitempty"`
	PrivacyLevel       *ScheduledEventPrivacyLevel                  `json:"privacy_level,omitempty"`
	ScheduledStartTime *time.Time                                   `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                                   `json:"scheduled_end_time,omitempty"`
	Description        *string                                      `json:"description,omitempty"`
	EntityType         *ScheduledEventEntityType                    `json:"entity_type,omitempty"`
	Status             *ScheduledEventStatus                        `json:"status,omitempty"`
	Image              *json.Nullable[Icon]                         `json:"image,omitempty"`
	RecurrenceRule     *json.Nullable[ScheduledEventRecurrenceRule] `json:"recurrence_rule,omitempty"`
}

type GuildScheduledEventUser struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	User                  User         `json:"user"`
	Member                *Member      `json:"member"`
}

// ScheduledEventPrivacyLevel the privacy level of the ScheduledEventPrivacyLevel (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level)
type ScheduledEventPrivacyLevel int

const (
	_ ScheduledEventPrivacyLevel = iota + 1
	ScheduledEventPrivacyLevelGuildOnly
)

// ScheduledEventStatus the status of the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status)
type ScheduledEventStatus int

const (
	ScheduledEventStatusScheduled ScheduledEventStatus = iota + 1
	ScheduledEventStatusActive
	ScheduledEventStatusCompleted
	ScheduledEventStatusCancelled
)

// ScheduledEventEntityType the type of the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types)
type ScheduledEventEntityType int

const (
	ScheduledEventEntityTypeStageInstance ScheduledEventEntityType = iota + 1
	ScheduledEventEntityTypeVoice
	ScheduledEventEntityTypeExternal
)

type ScheduledEventRecurrenceRule struct {
	Start      time.Time                              `json:"start"`
	End        *time.Time                             `json:"end"`
	Frequency  ScheduledEventRecurrenceRuleFrequency  `json:"frequency"`
	Interval   int                                    `json:"interval"`
	ByWeekday  []ScheduledEventRecurrenceRuleWeekday  `json:"by_weekday"`
	ByNWeekday []ScheduledEventRecurrenceRuleNWeekday `json:"by_n_weekday"`
	ByMonth    []ScheduledEventRecurrenceRuleMonth    `json:"by_month"`
	ByMonthDay []int                                  `json:"by_month_day"`
	ByYearDay  []int                                  `json:"by_year_day"`
	Count      *int                                   `json:"count"`
}

type ScheduledEventRecurrenceRuleFrequency int

const (
	ScheduledEventRecurrenceRuleFrequencyYearly ScheduledEventRecurrenceRuleFrequency = iota
	ScheduledEventRecurrenceRuleFrequencyMonthly
	ScheduledEventRecurrenceRuleFrequencyWeekly
	ScheduledEventRecurrenceRuleFrequencyDaily
)

type ScheduledEventRecurrenceRuleWeekday int

const (
	ScheduledEventRecurrenceRuleWeekdayMonday ScheduledEventRecurrenceRuleWeekday = iota
	ScheduledEventRecurrenceRuleWeekdayTuesday
	ScheduledEventRecurrenceRuleWeekdayWednesday
	ScheduledEventRecurrenceRuleWeekdayThursday
	ScheduledEventRecurrenceRuleWeekdayFriday
	ScheduledEventRecurrenceRuleWeekdaySaturday
	ScheduledEventRecurrenceRuleWeekdaySunday
)

type ScheduledEventRecurrenceRuleNWeekday struct {
	N   int                                 `json:"n"`
	Day ScheduledEventRecurrenceRuleWeekday `json:"day"`
}

type ScheduledEventRecurrenceRuleMonth int

const (
	ScheduledEventRecurrenceRuleMonthJanuary ScheduledEventRecurrenceRuleMonth = iota + 1
	ScheduledEventRecurrenceRuleMonthFebruary
	ScheduledEventRecurrenceRuleMonthMarch
	ScheduledEventRecurrenceRuleMonthApril
	ScheduledEventRecurrenceRuleMonthMay
	ScheduledEventRecurrenceRuleMonthJune
	ScheduledEventRecurrenceRuleMonthJuly
	ScheduledEventRecurrenceRuleMonthAugust
	ScheduledEventRecurrenceRuleMonthSeptember
	ScheduledEventRecurrenceRuleMonthOctober
	ScheduledEventRecurrenceRuleMonthNovember
	ScheduledEventRecurrenceRuleMonthDecember
)

// EntityMetaData additional metadata for the scheduled event (https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata)
type EntityMetaData struct {
	Location string `json:"location"`
}
