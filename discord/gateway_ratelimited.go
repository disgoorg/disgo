package discord

import (
	"github.com/disgoorg/snowflake/v2"
)

type RateLimitedMetadata interface {
	// ratelimitedmetadata is a marker to simulate unions.
	ratelimitedmetadata()
}

type RequestGuildMemberRateLimitMetadata struct {
	GuildID snowflake.ID `json:"guild_id"`
	Nonce   string       `json:"nonce"`
}

func (RequestGuildMemberRateLimitMetadata) ratelimitedmetadata() {}
