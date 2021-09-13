package discord

// Sticker is a sticker sent with a Message
type Sticker struct {
	ID          Snowflake         `json:"id"`
	PackID      *Snowflake        `json:"pack_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Tags        string            `json:"tags"`
	Type        StickerType       `json:"type"`
	FormatType  StickerFormatType `json:"format_type"`
	Available   *bool             `json:"available"`
	GuildID     *Snowflake        `json:"guild_id,omitempty"`
	User        *User             `json:"user,omitempty"`
	SortValue   *int              `json:"sort_value"`
}

type StickerType int

//goland:noinspection GoUnusedConst
const (
	StickerTypeStandard StickerType = iota + 1
	StickerTypeGuild
)

// StickerFormatType is the Format type of Sticker
type StickerFormatType int

// Constants for StickerFormatType
//goland:noinspection GoUnusedConst
const (
	StickerFormatTypePNG StickerFormatType = iota + 1
	StickerFormatTypeAPNG
	StickerFormatTypeLottie
)

type StickerCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags"`
	File        *File  `json:"-"`
}

// ToBody returns the MessageCreate ready for body
func (c *StickerCreate) ToBody() (interface{}, error) {
	if c.File != nil {
		return PayloadWithFiles(c, c.File)
	}
	return c, nil
}

type StickerUpdate struct {
	Name        *string         `json:"name,omitempty"`
	Description *OptionalString `json:"description,omitempty"`
	Tags        *string         `json:"tags,omitempty"`
}
