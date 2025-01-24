package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Entitlement struct {
	ID             snowflake.ID    `json:"id"`
	SkuID          snowflake.ID    `json:"sku_id"`
	ApplicationID  snowflake.ID    `json:"application_id"`
	UserID         *snowflake.ID   `json:"user_id"`
	PromotionID    *snowflake.ID   `json:"promotion_id"`
	Type           EntitlementType `json:"type"`
	Deleted        bool            `json:"deleted"`
	GiftCodeFlags  int             `json:"gift_code_flags"`
	Consumed       *bool           `json:"consumed"`
	StartsAt       *time.Time      `json:"starts_at"`
	EndsAt         *time.Time      `json:"ends_at"`
	GuildID        *snowflake.ID   `json:"guild_id"`
	SubscriptionID *snowflake.ID   `json:"subscription_id"`
}

type EntitlementType int

const (
	EntitlementTypePurchase EntitlementType = iota + 1
	EntitlementTypePremiumSubscription
	EntitlementTypeDeveloperGift
	EntitlementTypeTestModePurchase
	EntitlementTypeFreePurchase
	EntitlementTypeUserGift
	EntitlementTypePremiumPurchase
	EntitlementTypeApplicationSubscription
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
