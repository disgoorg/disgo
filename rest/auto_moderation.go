package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ AutoModeration = (*autoModerationImpl)(nil)

func NewAutoModeration(client Client) AutoModeration {
	return &autoModerationImpl{client: client}
}

type AutoModeration interface {
	GetAutoModerationRules(guildID snowflake.ID, opts ...RequestOpt) ([]discord.AutoModerationRule, error)
	GetAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, opts ...RequestOpt) (*discord.AutoModerationRule, error)
	CreateAutoModerationRule(guildID snowflake.ID, ruleCreate discord.AutoModerationRuleCreate, opts ...RequestOpt) (*discord.AutoModerationRule, error)
	UpdateAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, ruleUpdate discord.AutoModerationRuleUpdate, opts ...RequestOpt) (*discord.AutoModerationRule, error)
	DeleteAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, opts ...RequestOpt) error
}

type autoModerationImpl struct {
	client Client
}

func (s *autoModerationImpl) GetAutoModerationRules(guildID snowflake.ID, opts ...RequestOpt) (rules []discord.AutoModerationRule, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAutoModerationRules.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &rules, opts...)
	return
}

func (s *autoModerationImpl) GetAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, opts ...RequestOpt) (rule *discord.AutoModerationRule, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAutoModerationRule.Compile(nil, guildID, ruleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &rule, opts...)
	return
}

func (s *autoModerationImpl) CreateAutoModerationRule(guildID snowflake.ID, ruleCreate discord.AutoModerationRuleCreate, opts ...RequestOpt) (rule *discord.AutoModerationRule, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateAutoModerationRule.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, ruleCreate, &rule, opts...)
	return
}

func (s *autoModerationImpl) UpdateAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, ruleUpdate discord.AutoModerationRuleUpdate, opts ...RequestOpt) (rule *discord.AutoModerationRule, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateAutoModerationRule.Compile(nil, guildID, ruleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, ruleUpdate, &rule, opts...)
	return
}

func (s *autoModerationImpl) DeleteAutoModerationRule(guildID snowflake.ID, ruleID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteAutoModerationRule.Compile(nil, guildID, ruleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}
