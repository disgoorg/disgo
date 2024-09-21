package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var _ SKUs = (*skusImpl)(nil)

func NewSKUs(client Client) SKUs {
	return &skusImpl{client: client}
}

type SKUs interface {
	GetSKUs(applicationID snowflake.ID, opts ...RequestOpt) ([]discord.SKU, error)

	GetSKUSubscriptions(skuID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, userID snowflake.ID, opts ...RequestOpt) ([]discord.Subscription, error)
	GetSKUSubscriptionsPage(skuID snowflake.ID, userID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.Subscription]
	GetSKUSubscription(skuID snowflake.ID, subscriptionID snowflake.ID, opts ...RequestOpt) (*discord.Subscription, error)
}

type skusImpl struct {
	client Client
}

func (s *skusImpl) GetSKUs(applicationID snowflake.ID, opts ...RequestOpt) (skus []discord.SKU, err error) {
	err = s.client.Do(GetSKUs.Compile(nil, applicationID), nil, &skus, opts...)
	return
}

func (s *skusImpl) GetSKUSubscriptions(skuID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, userID snowflake.ID, opts ...RequestOpt) (subscriptions []discord.Subscription, err error) {
	values := discord.QueryValues{}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	if userID != 0 {
		values["user_id"] = userID
	}
	err = s.client.Do(GetSKUSubscriptions.Compile(values, skuID), nil, &subscriptions, opts...)
	return
}

func (s *skusImpl) GetSKUSubscriptionsPage(skuID snowflake.ID, userID snowflake.ID, startID snowflake.ID, limit int, opts ...RequestOpt) Page[discord.Subscription] {
	return Page[discord.Subscription]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) ([]discord.Subscription, error) {
			return s.GetSKUSubscriptions(skuID, before, after, limit, userID, opts...)
		},
		getIDFunc: func(subscription discord.Subscription) snowflake.ID {
			return subscription.ID
		},
		ID: startID,
	}
}

func (s *skusImpl) GetSKUSubscription(skuID snowflake.ID, subscriptionID snowflake.ID, opts ...RequestOpt) (subscription *discord.Subscription, err error) {
	err = s.client.Do(GetSKUSubscription.Compile(nil, skuID, subscriptionID), nil, &subscription, opts...)
	return
}
