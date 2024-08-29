package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Subscription struct {
	ID                 snowflake.ID       `json:"id"`
	UserID             snowflake.ID       `json:"user_id"`
	SkuIDs             []snowflake.ID     `json:"sku_ids"`
	EntitlementIDs     []snowflake.ID     `json:"entitlement_ids"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	Status             SubscriptionStatus `json:"status"`
	CanceledAt         *time.Time         `json:"canceled_at"`
	Country            *string            `json:"country"`
}

type SubscriptionStatus int

const (
	SubscriptionStatusActive SubscriptionStatus = iota
	SubscriptionStatusEnding
	SubscriptionStatusInactive
)
