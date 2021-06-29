package api

import "time"

var _ Thread = (*ChannelImpl)(nil)

type Thread interface {
	TextChannel
	MessageCount() int
	MemberCount() int
	ThreadMetadata() *ThreadMetadata
}

func (c *ChannelImpl) MessageCount() int {
	return c.MessageCount_
}

func (c *ChannelImpl) MemberCount() int {
	return c.MemberCount_
}

func (c *ChannelImpl) ThreadMetadata() *ThreadMetadata {
	return c.ThreadMetadata_
}

type ThreadMetadata struct {
	Archived            bool          `json:"archived"`
	ArchiveTimestamp    *time.Time    `json:"archive_timestamp"`
	ArchiverId          *Snowflake    `json:"archiver_id"`
	AutoArchiveDuration time.Duration `json:"auto_archive_duration"`
	Locked              bool          `json:"locked"`
}

type ThreadAutoArchiveDuration int

const (
	ThreadAutoArchiveDuration60    ThreadAutoArchiveDuration = 60
	ThreadAutoArchiveDuration1440  ThreadAutoArchiveDuration = 1440
	ThreadAutoArchiveDuration4320  ThreadAutoArchiveDuration = 4320
	ThreadAutoArchiveDuration10080 ThreadAutoArchiveDuration = 10080
)

type ThreadCreate struct {
	Name                string                    `json:"name"`
	AutoArchiveDuration ThreadAutoArchiveDuration `json:"auto_archive_duration"`
	Type                ChannelType               `json:"type"`
}
