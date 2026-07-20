package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// Attachment is used for files sent in a Message
type Attachment struct {
	ID                 snowflake.ID    `json:"id,omitempty"`
	Filename           string          `json:"filename,omitempty"`
	Title              *string         `json:"title,omitempty"`
	Description        *string         `json:"description,omitempty"`
	ContentType        *string         `json:"content_type,omitempty"`
	Size               int             `json:"size,omitempty"`
	URL                string          `json:"url,omitempty"`
	ProxyURL           string          `json:"proxy_url,omitempty"`
	Height             *int            `json:"height,omitempty"`
	Width              *int            `json:"width,omitempty"`
	Placeholder        string          `json:"placeholder,omitempty"`
	PlaceholderVersion int             `json:"placeholder_version,omitempty"`
	Ephemeral          bool            `json:"ephemeral,omitempty"`
	DurationSecs       *float64        `json:"duration_secs,omitempty"`
	Waveform           *string         `json:"waveform,omitempty"`
	Flags              AttachmentFlags `json:"flags"`
	ClipParticipants   []User          `json:"clip_participants,omitempty"`
	ClipCreatedAt      time.Time       `json:"clip_created_at,omitzero"`
	Application        *Application    `json:"application,omitempty"`
}

func (a Attachment) CreatedAt() time.Time {
	return a.ID.Time()
}

type AttachmentFlags int

const (
	AttachmentFlagIsClip AttachmentFlags = 1 << iota
	AttachmentFlagIsThumbnail
	AttachmentFlagIsRemix
	AttachmentFlagIsSpoiler
	_
	AttachmentFlagIsAnimated
	AttachmentFlagsNone AttachmentFlags = 0
)

// AttachmentUpdate is used in Message Create and Edit requests to set which Attachment(s) are in the message and provide their metadata.
type AttachmentUpdate interface {
	attachmentUpdate()
}

// AttachmentKeep is used to retain an existing Attachment when editing a Message.
//
// When editing, you must provide the ID of each file to keep, and you can optionally update the Description and IsSpoiler fields of existing attachments.
type AttachmentKeep struct {
	ID          snowflake.ID `json:"id"`
	Description *string      `json:"description,omitempty"`
	IsSpoiler   *bool        `json:"is_spoiler,omitempty"`
}

func (AttachmentKeep) attachmentUpdate() {}

// NewAttachmentKeep returns a new AttachmentKeep with the provided ID and no other fields set.
func NewAttachmentKeep(id snowflake.ID) AttachmentKeep {
	return AttachmentKeep{ID: id}
}

// WithDescription returns a new AttachmentKeep with the provided description.
func (a AttachmentKeep) WithDescription(description string) AttachmentKeep {
	a.Description = &description
	return a
}

// WithSpoiler returns a new AttachmentKeep with the provided spoiler setting.
func (a AttachmentKeep) WithSpoiler(spoiler bool) AttachmentKeep {
	a.IsSpoiler = &spoiler
	return a
}

// AttachmentCreate is used to describe the metadata of a newly uploaded file when creating or editing a Message.
type AttachmentCreate struct {
	ID          int    `json:"id"`
	Description string `json:"description,omitempty"`
	IsSpoiler   *bool  `json:"is_spoiler,omitempty"`
}

func (AttachmentCreate) attachmentUpdate() {}

// WithDescription returns a new AttachmentCreate with the provided description.
func (a AttachmentCreate) WithDescription(description string) AttachmentCreate {
	a.Description = description
	return a
}

// WithSpoiler returns a new AttachmentCreate with the provided spoiler setting.
func (a AttachmentCreate) WithSpoiler(spoiler bool) AttachmentCreate {
	a.IsSpoiler = &spoiler
	return a
}
