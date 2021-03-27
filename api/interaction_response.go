package api

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong InteractionResponseType = iota + 1
	_
	_
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
)

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

// The InteractionResponseData is used to specify the message_events options when creating an InteractionResponse
type InteractionResponseData struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
}

// InteractionResponseBuilder allows you to create an InteractionResponse with ease
type InteractionResponseBuilder struct {
	InteractionResponse
}

// NewInteractionResponseBuilder returns a new InteractionResponseBuilder
func NewInteractionResponseBuilder() *InteractionResponseBuilder {
	return &InteractionResponseBuilder{
		InteractionResponse{
			Type: InteractionResponseTypeChannelMessageWithSource,
			Data: &InteractionResponseData{
				AllowedMentions: &DefaultInteractionAllowedMentions,
			},
		},
	}
}

// SetType sets the InteractionResponseType of the InteractionResponse
func (b *InteractionResponseBuilder) SetType(responseType InteractionResponseType) *InteractionResponseBuilder {
	b.Type = responseType
	return b
}

// SetData sets the InteractionResponseData of the InteractionResponse
func (b *InteractionResponseBuilder) SetData(data *InteractionResponseData) *InteractionResponseBuilder {
	b.Data = data
	return b
}

// SetTTS sets if the InteractionResponse is a tts message
func (b *InteractionResponseBuilder) SetTTS(tts bool) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.TTS = tts
	return b
}

// SetContent sets the content of the InteractionResponse
func (b *InteractionResponseBuilder) SetContent(content string) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.Content = content
	return b
}

// SetEmbeds sets the embeds of the InteractionResponse
func (b *InteractionResponseBuilder) SetEmbeds(embeds ...*Embed) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the InteractionResponse
func (b *InteractionResponseBuilder) AddEmbeds(embeds ...*Embed) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.Embeds = append(b.Data.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the InteractionResponse
func (b *InteractionResponseBuilder) ClearEmbeds() *InteractionResponseBuilder {
	if b.Data != nil && b.Data.Embeds != nil {
		b.Data.Embeds = []*Embed{}
	}
	return b
}

// RemoveEmbed removes an embed from the InteractionResponse
func (b *InteractionResponseBuilder) RemoveEmbed(index int) *InteractionResponseBuilder {
	if b.Data != nil && len(b.Data.Embeds) > index {
		b.Data.Embeds = append(b.Data.Embeds[:index], b.Data.Embeds[index+1:]...)
	}
	return b
}

// SetAllowedMentions sets the allowed mentions of the InteractionResponse
func (b *InteractionResponseBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the InteractionResponse to nothing
func (b *InteractionResponseBuilder) SetAllowedMentionsEmpty() *InteractionResponseBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the InteractionResponse
func (b *InteractionResponseBuilder) SetFlags(flags MessageFlags) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	b.Data.Flags = flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *InteractionResponseBuilder) SetEphemeral(ephemeral bool) *InteractionResponseBuilder {
	if b.Data == nil {
		b.Data = &InteractionResponseData{}
	}
	if ephemeral {
		b.Data.Flags |= MessageFlagEphemeral
	} else {
		b.Data.Flags &^= MessageFlagEphemeral
	}
	return b
}

// Build returns your built InteractionResponse
func (b *InteractionResponseBuilder) Build() *InteractionResponse {
	return &b.InteractionResponse
}

// FollowupMessage is used to add additional messages to an Interaction after you've responded initially
type FollowupMessage struct {
	// Todo: fill this
}
