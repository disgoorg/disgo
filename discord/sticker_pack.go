package discord

import "github.com/DisgoOrg/snowflake"

type StickerPack struct {
	ID             snowflake.Snowflake `json:"id"`
	Stickers       []Sticker           `json:"stickers"`
	Name           string              `json:"name"`
	SkuID          snowflake.Snowflake `json:"sku_id"`
	CoverStickerID snowflake.Snowflake `json:"cover_sticker_id"`
	Description    string              `json:"description"`
	BannerAssetID  snowflake.Snowflake `json:"banner_asset_id"`
}

type StickerPacks struct {
	StickerPacks []StickerPack `json:"sticker_packs"`
}
