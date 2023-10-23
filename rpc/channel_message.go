package rpc

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type ChannelMessage struct {
	ID              snowflake.ID   `json:"id"`
	Blocked         bool           `json:"blocked"`
	Bot             bool           `json:"bot"`
	Content         string         `json:"content"`
	ContentParsed   []interface{}  `json:"content_parsed"`
	Nick            string         `json:"nick"`
	EditedTimestamp time.Time      `json:"edited_timestamp"`
	Timestamp       time.Time      `json:"timestamp"`
	TTS             bool           `json:"tts"`
	Mentions        []snowflake.ID `json:"mentions"`
	MentionEveryone bool           `json:"mention_everyone"`
	MentionRoles    []snowflake.ID `json:"mention_roles"`
	Embeds          []Embed        `json:"embeds"`
	Attachments     []Attachment   `json:"attachments"`
	Author          Author         `json:"author"`
	Pinned          bool           `json:"pinned"`
	Type            int            `json:"type"`
}

func (m *ChannelMessage) UnmarshalJSON(data []byte) error {
	type message ChannelMessage
	var v struct {
		message
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*m = ChannelMessage(v.message)
	return nil
}

type AvatarDecorationData struct {
	Asset string `json:"asset"`
	SkuID string `json:"sku_id"`
}

type Author struct {
	ID                   snowflake.ID         `json:"id"`
	Username             string               `json:"username"`
	Discriminator        string               `json:"discriminator"`
	GlobalName           string               `json:"global_name"`
	Avatar               string               `json:"avatar"`
	AvatarDecorationData AvatarDecorationData `json:"avatar_decoration_data"`
	Bot                  bool                 `json:"bot"`
	Flags                int                  `json:"flags"`
	PremiumType          discord.PremiumType  `json:"premium_type"`
}
