package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// Interaction represents a generic Interaction received from discord
type Interaction interface {
	discord.Interaction
}

type InteractionFields struct {
	discord.InteractionFields
	Bot             *Bot
	User            *User
	Member          *Member
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
}

// ApplicationCommandInteraction represents a generic ApplicationCommand Interaction received from discord
type ApplicationCommandInteraction interface {
	discord.ApplicationCommandInteraction
}

// ComponentInteraction represents a generic discord.Component Interaction received from discord
type ComponentInteraction interface {
	discord.ComponentInteraction
}

func respond(fields *InteractionFields, callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	if fields.Acknowledged {
		return discord.ErrInteractionAlreadyReplied
	}
	fields.Acknowledged = true

	response := discord.InteractionResponse{
		Type: callbackType,
		Data: callbackData,
	}

	if fields.ResponseChannel != nil {
		fields.ResponseChannel <- response
		return nil
	}

	return fields.Bot.RestServices.InteractionService().CreateInteractionResponse(fields.ID, fields.Token, response, opts...)
}

func deferCreate(fields *InteractionFields, ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return respond(fields, discord.InteractionCallbackTypeDeferredChannelMessageWithSource, data, opts...)
}

func create(fields *InteractionFields, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return respond(fields, discord.InteractionCallbackTypeChannelMessageWithSource, messageCreate, opts...)
}

func getOriginal(fields *InteractionFields, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices.InteractionService().GetInteractionResponse(fields.ApplicationID, fields.Token, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func updateOriginal(fields *InteractionFields, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices.InteractionService().UpdateInteractionResponse(fields.ApplicationID, fields.Token, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func deleteOriginal(fields *InteractionFields, opts ...rest.RequestOpt) error {
	return fields.Bot.RestServices.InteractionService().DeleteInteractionResponse(fields.ApplicationID, fields.Token, opts...)
}

func deferUpdate(fields *InteractionFields, opts ...rest.RequestOpt) error {
	return respond(fields, discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

func update(fields *InteractionFields, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return respond(fields, discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func updateComponent(fields *InteractionFields, message *Message, customID string, component discord.Component, opts ...rest.RequestOpt) error {
	actionRows := message.ActionRows()
	for _, actionRow := range actionRows {
		actionRow = actionRow.SetComponent(customID, component)
	}

	return update(fields, NewMessageUpdateBuilder().SetActionRows(actionRows...).Build(), opts...)
}

func result(fields *InteractionFields, choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return respond(fields, discord.InteractionCallbackTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

func resultMapString(fields *InteractionFields, resultMap map[string]string, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceString{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, choices, opts...)
}

func resultMapInt(fields *InteractionFields, resultMap map[string]int, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceInt{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, choices, opts...)
}

func resultMapFloat(fields *InteractionFields, resultMap map[string]float64, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceFloat{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, choices, opts...)
}

func createFollowup(fields *InteractionFields, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices.InteractionService().CreateFollowupMessage(fields.ApplicationID, fields.Token, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func updateFollowup(fields *InteractionFields, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices.InteractionService().UpdateFollowupMessage(fields.ApplicationID, fields.Token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func deleteFollowup(fields *InteractionFields, messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return fields.Bot.RestServices.InteractionService().DeleteFollowupMessage(fields.ApplicationID, fields.Token, messageID, opts...)
}
