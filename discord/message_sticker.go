package discord

type MessageSticker struct {
	ID         Snowflake         `json:"id"`
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}
