package rpc

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
)

// The EmbedResource of an Embed.Image/Embed.Thumbnail/Embed.Video
type EmbedResource struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxyURL,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type Embed struct {
	discord.Embed
}

func (e *Embed) UnmarshalJSON(data []byte) error {
	type embed Embed
	var v struct {
		embed
		RawTitle       string         `json:"rawTitle"`
		RawDescription string         `json:"rawDescription"`
		Image          *EmbedResource `json:"image,omitempty"`
		Thumbnail      *EmbedResource `json:"thumbnail,omitempty"`
		Video          *EmbedResource `json:"video,omitempty"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = Embed(v.embed)
	e.Title = v.RawTitle
	e.Description = v.RawDescription
	if v.Image != nil {
		e.Image = &discord.EmbedResource{
			URL:      v.Image.URL,
			ProxyURL: v.Image.ProxyURL,
			Height:   v.Image.Height,
			Width:    v.Image.Width,
		}
	}
	if v.Thumbnail != nil {
		e.Thumbnail = &discord.EmbedResource{
			URL:      v.Thumbnail.URL,
			ProxyURL: v.Thumbnail.ProxyURL,
			Height:   v.Thumbnail.Height,
			Width:    v.Thumbnail.Width,
		}
	}
	if v.Video != nil {
		e.Video = &discord.EmbedResource{
			URL:      v.Video.URL,
			ProxyURL: v.Video.ProxyURL,
			Height:   v.Video.Height,
			Width:    v.Video.Width,
		}
	}
	return nil
}
