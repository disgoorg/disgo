package gateway

import (
	"github.com/disgoorg/snowflake/v2"
)

// RateLimitedMetadata is an interface for all ratelimited metadatas.
// [RateLimitedMetadataRequestGuildMembers]
// [RateLimitedMetadataUnknown]
type RateLimitedMetadata interface {
	// ratelimitedmetadata is a marker to simulate unions.
	ratelimitedmetadata()
}

type RateLimitedMetadataRequestGuildMembers struct {
	GuildID snowflake.ID `json:"guild_id"`
	Nonce   string       `json:"nonce"`
}

func (RateLimitedMetadataRequestGuildMembers) ratelimitedmetadata() {}

type RateLimitedMetadataUnknown struct{}

func (RateLimitedMetadataUnknown) ratelimitedmetadata() {}
