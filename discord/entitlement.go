package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Entitlement struct {
	ID             snowflake.ID    `json:"id"`
	SkuID          snowflake.ID    `json:"sku_id"`
	UserID         *snowflake.ID   `json:"user_id"`
	GuildID        *snowflake.ID   `json:"guild_id"`
	ApplicationID  snowflake.ID    `json:"application_id"`
	Type           EntitlementType `json:"type"`
	Consumed       bool            `json:"consumed"`
	StartsAt       *time.Time      `json:"starts_at"`
	EndsAt         *time.Time      `json:"ends_at"`
	PromotionID    *snowflake.ID   `json:"promotion_id"`
	Deleted        bool            `json:"deleted"`
	GiftCodeFlags  int             `json:"gift_code_flags"`
	SubscriptionID *snowflake.ID   `json:"subscription_id"`
}

type EntitlementType int

const (
	EntitlementTypeApplicationSubscription EntitlementType = 8
)

type TestEntitlementCreate struct {
	SkuID     snowflake.ID         `json:"sku_id"`
	OwnerID   snowflake.ID         `json:"owner_id"`
	OwnerType EntitlementOwnerType `json:"owner_type"`
}

type EntitlementOwnerType int

const (
	EntitlementOwnerTypeGuild EntitlementOwnerType = iota + 1
	EntitlementOwnerTypeUser
)
