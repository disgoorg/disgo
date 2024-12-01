package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/slicehelper"
	"github.com/disgoorg/snowflake/v2"
)

// QueryParams serves as a generic interface for implementations of rest endpoint query parameters.
type QueryParams interface {
	// ToQueryValues transforms fields from the QueryParams interface implementations into discord.QueryValues.
	ToQueryValues() discord.QueryValues
}

// QueryParamsGetEntitlements holds query parameters for Applications.GetEntitlements (https://discord.com/developers/docs/resources/entitlement#list-entitlements)
type QueryParamsGetEntitlements struct {
	UserID         snowflake.ID
	SkuIDs         []snowflake.ID
	Before         int
	After          int
	Limit          int
	GuildID        snowflake.ID
	ExcludeEnded   bool
	ExcludeDeleted bool
}

func (p QueryParamsGetEntitlements) ToQueryValues() discord.QueryValues {
	queryValues := discord.QueryValues{
		"exclude_ended":   p.ExcludeEnded,
		"exclude_deleted": p.ExcludeDeleted,
		"sku_ids":         slicehelper.JoinSnowflakes(p.SkuIDs),
	}
	if p.UserID != 0 {
		queryValues["user_id"] = p.UserID
	}
	if p.Before != 0 {
		queryValues["before"] = p.Before
	}
	if p.After != 0 {
		queryValues["after"] = p.After
	}
	if p.Limit != 0 {
		queryValues["limit"] = p.Limit
	}
	if p.GuildID != 0 {
		queryValues["guild_id"] = p.GuildID
	}
	return queryValues
}
