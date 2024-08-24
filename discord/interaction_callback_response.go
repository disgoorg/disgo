package discord

import "github.com/disgoorg/snowflake/v2"

type InteractionCallbackResponse struct {
	Interaction CallbackInteraction `json:"interaction"`
	Resource    *CallbackResource   `json:"resource"`
}

type CallbackInteraction struct {
	ID                       snowflake.ID    `json:"id"`
	Type                     InteractionType `json:"type"`
	ResponseMessageID        *snowflake.ID   `json:"response_message_id"`
	ResponseMessageLoading   *bool           `json:"response_message_loading"`
	ResponseMessageEphemeral *bool           `json:"response_message_ephemeral"`
}

type CallbackResource struct {
	Type    InteractionResponseType `json:"type"`
	Message *Message                `json:"message"`
}
