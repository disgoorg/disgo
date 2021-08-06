package api

type StickerPack struct {
	ID             Snowflake         `json:"id"`
	Stickers       []*MessageSticker `json:"stickers"`
	Name           string            `json:"name"`
	SkuID          Snowflake         `json:"sku_id"`
	CoverStickerID Snowflake         `json:"cover_sticker_id,omitempty"`
	Description    string            `json:"description"`
	BannerAssetID  Snowflake         `json:"banner_asset_id"`
}
