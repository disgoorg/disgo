package discord

import (
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

// Sticker is a sticker sent with a Message
type Sticker struct {
	ID          snowflake.Snowflake  `json:"id"`
	PackID      *snowflake.Snowflake `json:"pack_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Tags        string               `json:"tags"`
	Type        StickerType          `json:"type"`
	FormatType  StickerFormatType    `json:"format_type"`
	Available   *bool                `json:"available"`
	GuildID     *snowflake.Snowflake `json:"guild_id,omitempty"`
	User        *User                `json:"user,omitempty"`
	SortValue   *int                 `json:"sort_value"`
}

func (s Sticker) URL(opts ...CDNOpt) string {
	format := route.PNG
	if s.FormatType == StickerFormatTypeLottie {
		format = route.Lottie
	}
	if avatar := formatAssetURL(route.CustomSticker, append(opts, WithFormat(format)), s.ID); avatar != nil {
		return *avatar
	}
	return ""
}

type StickerType int

const (
	StickerTypeStandard StickerType = iota + 1
	StickerTypeGuild
)

// StickerFormatType is the Format type of Sticker
type StickerFormatType int

// Constants for StickerFormatType
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
func (c *StickerCreate) ToBody() (any, error) {
	if c.File != nil {
		return PayloadWithFiles(c, c.File)
	}
	return c, nil
}

type StickerUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

type StickerPack struct {
	ID             snowflake.Snowflake `json:"id"`
	Stickers       []Sticker           `json:"stickers"`
	Name           string              `json:"name"`
	SkuID          snowflake.Snowflake `json:"sku_id"`
	CoverStickerID snowflake.Snowflake `json:"cover_sticker_id"`
	Description    string              `json:"description"`
	BannerAssetID  snowflake.Snowflake `json:"banner_asset_id"`
}

func (p StickerPack) BannerURL(opts ...CDNOpt) *string {
	return formatAssetURL(route.StickerPackBanner, opts, p.BannerAssetID)
}

type StickerPacks struct {
	StickerPacks []StickerPack `json:"sticker_packs"`
}
