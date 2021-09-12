package discord

// StickerFormatType is the Format type of Sticker
type StickerFormatType int

// Constants for StickerFormatType
//goland:noinspection GoUnusedConst
const (
	StickerFormatPNG StickerFormatType = iota + 1
	StickerFormatAPNG
	StickerFormatLottie
)
