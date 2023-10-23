package rpc

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

// The EmbedResource of an Embed.Image/Embed.Thumbnail/Embed.Video
type EmbedResource struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxyURL,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type Embed struct {
	ID             string                 `json:"id"`   // eg. "embed_N"
	Type           discord.EmbedType      `json:"type"` // Only seen "rich" so far
	RawTitle       string                 `json:"rawTitle"`
	RawDescription string                 `json:"rawDescription"`
	Color          string                 `json:"color"` // CSS color, why discord why.
	Fields         []discord.EmbedField   `json:"fields"`
	URL            string                 `json:"url,omitempty"`
	Timestamp      *time.Time             `json:"timestamp,omitempty"`
	Footer         *discord.EmbedFooter   `json:"footer,omitempty"`
	Image          *EmbedResource         `json:"image,omitempty"`
	Thumbnail      *EmbedResource         `json:"thumbnail,omitempty"`
	Video          *EmbedResource         `json:"video,omitempty"`
	Provider       *discord.EmbedProvider `json:"provider,omitempty"`
	Author         *discord.EmbedAuthor   `json:"author,omitempty"`
}

type Attachment struct {
	*discord.Attachment
	Spoiler bool `json:"spoiler,omitempty"`
}
