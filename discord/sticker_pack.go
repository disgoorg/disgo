package discord

type StickerPack struct {
	ID             Snowflake `json:"id"`
	Stickers       []Sticker `json:"stickers"`
	Name           string    `json:"name"`
	SkuID          Snowflake `json:"sku_id"`
	CoverStickerID Snowflake `json:"cover_sticker_id"`
	Description    string    `json:"description"`
	BannerAssetID  Snowflake `json:"banner_asset_id"`
}

type StickerPacks struct {
	StickerPacks []StickerPack `json:"sticker_packs"`
}
