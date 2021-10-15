package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type InteractionData struct {
	Bot             *Bot
	User            *User
	Member          *Member
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
}

// Interaction represents a generic Interaction received from discord
type Interaction interface {
	discord.Interaction
}

type ApplicationCommandInteraction interface {
	discord.ApplicationCommandInteraction
}

type ComponentInteraction interface {
	discord.ComponentInteraction
}

type RespondInteraction struct {
	InteractionData
	id    discord.Snowflake
	token string
}

// Respond responds to the Interaction with the provided discord.InteractionResponse
func (i *RespondInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	response := discord.InteractionResponse{
		Type: callbackType,
		Data: callbackData,
	}
	if i.Acknowledged {
		return discord.ErrInteractionAlreadyReplied
	}
	i.Acknowledged = true

	if i.ResponseChannel != nil {
		i.ResponseChannel <- response
		return nil
	}

	return i.Bot.RestServices.InteractionService().CreateInteractionResponse(i.id, i.token, response, opts...)
}

type CreateInteraction struct {
	RespondInteraction
	applicationID discord.Snowflake
}

// DeferCreate replies to the Interaction with discord.InteractionCallbackTypeDeferredChannelMessageWithSource and shows a loading state
func (i *CreateInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(discord.InteractionCallbackTypeDeferredChannelMessageWithSource, data, opts...)
}

// Create replies to the Interaction with discord.InteractionCallbackTypeChannelMessageWithSource & discord.MessageCreate
func (i *CreateInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeChannelMessageWithSource, messageCreate, opts...)
}

// GetOriginal gets the original discord.InteractionResponse
func (i *CreateInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().GetInteractionResponse(i.applicationID, i.token, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateOriginal edits the original discord.InteractionResponse
func (i *CreateInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateInteractionResponse(i.applicationID, i.token, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (i *CreateInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.InteractionService().DeleteInteractionResponse(i.applicationID, i.token, opts...)
}

type UpdateInteraction struct {
	RespondInteraction
	message  *Message
	customID string
}

// DeferUpdate replies to the ComponentInteraction with discord.InteractionCallbackTypeDeferredUpdateMessage and cancels the loading state
func (i *UpdateInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

// Update replies to the ComponentInteraction with discord.InteractionCallbackTypeUpdateMessage & discord.MessageUpdate which edits the original Message
func (i *UpdateInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (i *UpdateInteraction) UpdateComponent(component discord.Component, opts ...rest.RequestOpt) error {
	actionRows := i.message.ActionRows()
	for _, actionRow := range actionRows {
		actionRow = actionRow.SetComponent(i.customID, component)
	}

	return i.Update(NewMessageUpdateBuilder().SetActionRows(actionRows...).Build(), opts...)
}

type ResultInteraction struct {
	RespondInteraction
	applicationID discord.Snowflake
}

func (i *ResultInteraction) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

func (i *ResultInteraction) ResultMap(resultMap map[string]string, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoice{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return i.Result(choices, opts...)
}

type FollowupInteraction struct {
	InteractionData
	id            discord.Snowflake
	token         string
	applicationID discord.Snowflake
}

// CreateFollowup is used to send a discord.MessageCreate to an Interaction
func (i *FollowupInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().CreateFollowupMessage(i.applicationID, i.token, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateFollowup is used to edit a Message from an Interaction
func (i *FollowupInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateFollowupMessage(i.applicationID, i.token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *FollowupInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.InteractionService().DeleteFollowupMessage(i.applicationID, i.token, messageID, opts...)
}
