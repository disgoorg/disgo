package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type SKU struct {
	ID             snowflake.ID  `json:"id"`
	Type           SKUType       `json:"type"`
	ApplicationID  snowflake.ID  `json:"application_id"`
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	DependentSkuID *snowflake.ID `json:"dependent_sku_id"`
	AccessType     int           `json:"access_type"`
	Features       []string      `json:"features"`
	ReleaseDate    *time.Time    `json:"release_date"`
	Premium        bool          `json:"premium"`
	Flags          SKUFlags      `json:"flags"`
	ShowAgeGate    bool          `json:"show_age_gate"`
}

type SKUType int

const (
	SKUTypeDurable SKUType = iota + 2
	SKUTypeConsumable
	_
	SKUTypeSubscription
	SKUTypeSubscriptionGroup
)

type SKUFlags int

const (
	SKUFlagAvailable SKUFlags = 1 << (iota + 2)
	_
	_
	_
	_
	SKUFlagGuildSubscription
	SKUFlagUserSubscription
)
