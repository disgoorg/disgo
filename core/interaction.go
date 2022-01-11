package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type InteractionFilter func(interaction Interaction) bool

// Interaction represents a generic Interaction received from discord
type Interaction interface {
	discord.Interaction
}

type InteractionFields struct {
	Bot             Bot
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
}

type ApplicationCommandInteractionFilter func(interaction Interaction) bool

// ApplicationCommandInteraction represents a generic ApplicationCommand Interaction received from discord
type ApplicationCommandInteraction interface {
	discord.ApplicationCommandInteraction
}

type ComponentInteractionFilter func(interaction Interaction) bool

// ComponentInteraction represents a generic discord.Component Interaction received from discord
type ComponentInteraction interface {
	discord.ComponentInteraction
}

func commandPath(commandName string, subCommandName *string, subCommandGroupName *string) string {
	path := commandName
	if name := subCommandName; name != nil {
		path += "/" + *name
	}
	if name := subCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

func respond(fields *InteractionFields, id discord.Snowflake, token string, callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
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

	return fields.Bot.RestServices().InteractionService().CreateInteractionResponse(id, token, response, opts...)
}

func deferCreate(fields *InteractionFields, id discord.Snowflake, token string, ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return respond(fields, id, token, discord.InteractionCallbackTypeDeferredChannelMessageWithSource, data, opts...)
}

func create(fields *InteractionFields, id discord.Snowflake, token string, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return respond(fields, id, token, discord.InteractionCallbackTypeChannelMessageWithSource, messageCreate, opts...)
}

func getOriginal(fields *InteractionFields, applicationID discord.Snowflake, token string, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices().InteractionService().GetInteractionResponse(applicationID, token, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

func updateOriginal(fields *InteractionFields, applicationID discord.Snowflake, token string, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices().InteractionService().UpdateInteractionResponse(applicationID, token, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

func deleteOriginal(fields *InteractionFields, applicationID discord.Snowflake, token string, opts ...rest.RequestOpt) error {
	return fields.Bot.RestServices().InteractionService().DeleteInteractionResponse(applicationID, token, opts...)
}

func deferUpdate(fields *InteractionFields, applicationID discord.Snowflake, token string, opts ...rest.RequestOpt) error {
	return respond(fields, applicationID, token, discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

func update(fields *InteractionFields, applicationID discord.Snowflake, token string, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return respond(fields, applicationID, token, discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func updateComponent(fields *InteractionFields, applicationID discord.Snowflake, token string, message *Message, customID discord.CustomID, component discord.InteractiveComponent, opts ...rest.RequestOpt) error {
	containerComponents := make([]discord.ContainerComponent, len(message.Components))
	for i := range message.Components {
		switch container := containerComponents[i].(type) {
		case discord.ActionRowComponent:
			containerComponents[i] = container.UpdateComponent(customID, component)

		default:
			containerComponents[i] = container
			continue
		}
	}

	return update(fields, applicationID, token, discord.NewMessageUpdateBuilder().SetContainerComponents(containerComponents...).Build(), opts...)
}

func result(fields *InteractionFields, applicationID discord.Snowflake, token string, choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return respond(fields, applicationID, token, discord.InteractionCallbackTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

func resultMapString(fields *InteractionFields, applicationID discord.Snowflake, token string, resultMap map[string]string, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceString{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, applicationID, token, choices, opts...)
}

func resultMapInt(fields *InteractionFields, applicationID discord.Snowflake, token string, resultMap map[string]int, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceInt{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, applicationID, token, choices, opts...)
}

func resultMapFloat(fields *InteractionFields, applicationID discord.Snowflake, token string, resultMap map[string]float64, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceFloat{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return result(fields, applicationID, token, choices, opts...)
}

func getFollowup(fields *InteractionFields, applicationID discord.Snowflake, token string, messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices().InteractionService().GetFollowupMessage(applicationID, token, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

func createFollowup(fields *InteractionFields, applicationID discord.Snowflake, token string, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices().InteractionService().CreateFollowupMessage(applicationID, token, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

func updateFollowup(fields *InteractionFields, applicationID discord.Snowflake, token string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := fields.Bot.RestServices().InteractionService().UpdateFollowupMessage(applicationID, token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return fields.Bot.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

func deleteFollowup(fields *InteractionFields, applicationID discord.Snowflake, token string, messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return fields.Bot.RestServices().InteractionService().DeleteFollowupMessage(applicationID, token, messageID, opts...)
}

func channel(fields *InteractionFields, channelID discord.Snowflake) MessageChannel {
	if ch := fields.Bot.Caches().Channels().Get(channelID); ch != nil {
		return ch.(MessageChannel)
	}
	return nil
}

func guild(fields *InteractionFields, guildID *discord.Snowflake) *Guild {
	if guildID == nil {
		return nil
	}
	return fields.Bot.Caches().Guilds().Get(*guildID)
}
