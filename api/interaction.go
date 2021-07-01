package api

import (
	"encoding/json"
	"errors"

	"github.com/DisgoOrg/restclient"
)

// InteractionType is the type of Interaction
type InteractionType int

// Supported InteractionType(s)
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeCommand
	InteractionTypeComponent
)

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong InteractionResponseType = iota + 1
	_
	_
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
	InteractionResponseTypeDeferredUpdateMessage
	InteractionResponseTypeUpdateMessage
)

// Interaction holds the general parameters of each Interaction
type Interaction struct {
	Disgo           Disgo
	ResponseChannel chan InteractionResponse
	Replied         bool
	ID              Snowflake       `json:"id"`
	Type            InteractionType `json:"type"`
	GuildID         *Snowflake      `json:"guild_id,omitempty"`
	ChannelID       *Snowflake      `json:"channel_id,omitempty"`
	Member          *Member         `json:"member,omitempty"`
	User            *User           `json:"User,omitempty"`
	Token           string          `json:"token"`
	Version         int             `json:"version"`
}

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType `json:"type"`
	Data interface{}             `json:"data,omitempty"`
}

// ToBody returns the InteractionResponse ready for body
func (r InteractionResponse) ToBody() (interface{}, error) {
	if r.Data == nil {
		return r, nil
	}
	switch v := r.Data.(type) {
	case MessageCreate:
		if len(v.Files) > 0 {
			return restclient.PayloadWithFiles(r, v.Files...)
		}
	case MessageUpdate:
		if len(v.Files) > 0 {
			return restclient.PayloadWithFiles(r, v.Files...)
		}
	}
	return r, nil
}

// Respond responds to the api.Interaction with the provided api.InteractionResponse
func (i *Interaction) Respond(responseType InteractionResponseType, data interface{}) restclient.RestError {
	response := InteractionResponse{
		Type: responseType,
		Data: data,
	}
	if i.Replied {
		return restclient.NewError(nil, errors.New("you already replied to this interaction"))
	}
	i.Replied = true

	if i.FromWebhook() {
		i.ResponseChannel <- response
		return nil
	}

	return i.Disgo.RestClient().CreateInteractionResponse(i.ID, i.Token, response)
}

// DeferReply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (i *Interaction) DeferReply(ephemeral bool) restclient.RestError {
	var messageCreate interface{}
	if ephemeral {
		messageCreate = MessageCreate{Flags: MessageFlagEphemeral}
	}
	return i.Respond(InteractionResponseTypeDeferredChannelMessageWithSource, messageCreate)
}

// Reply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (i *Interaction) Reply(messageCreate MessageCreate) restclient.RestError {
	return i.Respond(InteractionResponseTypeChannelMessageWithSource, messageCreate)
}

// EditOriginal edits the original api.InteractionResponse
func (i *Interaction) EditOriginal(messageUpdate MessageUpdate) (*Message, restclient.RestError) {
	return i.Disgo.RestClient().UpdateInteractionResponse(i.Disgo.ApplicationID(), i.Token, messageUpdate)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (i *Interaction) DeleteOriginal() restclient.RestError {
	return i.Disgo.RestClient().DeleteInteractionResponse(i.Disgo.ApplicationID(), i.Token)
}

// SendFollowup used to send a api.MessageCreate to an api.Interaction
func (i *Interaction) SendFollowup(messageCreate MessageCreate) (*Message, restclient.RestError) {
	return i.Disgo.RestClient().CreateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageCreate)
}

// EditFollowup used to edit a api.Message from an api.Interaction
func (i *Interaction) EditFollowup(messageID Snowflake, messageUpdate MessageUpdate) (*Message, restclient.RestError) {
	return i.Disgo.RestClient().UpdateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID, messageUpdate)
}

// DeleteFollowup used to delete a api.Message from an api.Interaction
func (i *Interaction) DeleteFollowup(messageID Snowflake) restclient.RestError {
	return i.Disgo.RestClient().DeleteFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID)
}

// FromWebhook returns is the Interaction was made via http
func (i *Interaction) FromWebhook() bool {
	return i.ResponseChannel != nil
}

// Guild returns the api.Guild from the api.Cache
func (i *Interaction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Disgo.Cache().Guild(*i.GuildID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (i *Interaction) DMChannel() DMChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().DMChannel(*i.ChannelID)
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (i *Interaction) MessageChannel() MessageChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().MessageChannel(*i.ChannelID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (i *Interaction) TextChannel() TextChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().TextChannel(*i.ChannelID)
}

// GuildChannel returns the api.GuildChannel from the api.Cache
func (i *Interaction) GuildChannel() GuildChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().GuildChannel(*i.ChannelID)
}

// FullInteraction is used for easier unmarshalling of different Interaction(s)
type FullInteraction struct {
	ID        Snowflake       `json:"id"`
	Type      InteractionType `json:"type"`
	GuildID   *Snowflake      `json:"guild_id,omitempty"`
	ChannelID *Snowflake      `json:"channel_id,omitempty"`
	Message   *Message        `json:"message,omitempty"`
	Member    *Member         `json:"member,omitempty"`
	User      *User           `json:"User,omitempty"`
	Token     string          `json:"token"`
	Version   int             `json:"version"`
	Data      json.RawMessage `json:"data,omitempty"`
}
