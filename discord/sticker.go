package discord

// Sticker is a sticker sent with a Message
type Sticker struct {
	ID          Snowflake         `json:"id"`
	PackID      Snowflake         `json:"pack_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        *string           `json:"tags"`
	FormatType  StickerFormatType `json:"format_type"`
}
