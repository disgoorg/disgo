package discord

import "github.com/disgoorg/snowflake/v2"

type InteractionCallbackResponse struct {
	Interaction InteractionCallback          `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource"`
}

type InteractionCallback struct {
	ID                       snowflake.ID    `json:"id"`
	Type                     InteractionType `json:"type"`
	ActivityInstanceID       string          `json:"activity_instance_id"`
	ResponseMessageID        snowflake.ID    `json:"response_message_id"`
	ResponseMessageLoading   bool            `json:"response_message_loading"`
	ResponseMessageEphemeral bool            `json:"response_message_ephemeral"`
}

type InteractionCallbackResource struct {
	Type             InteractionResponseType              `json:"type"`
	ActivityInstance *InteractionCallbackActivityInstance `json:"activity_instance"`
	Message          *Message                             `json:"message"`
}

type InteractionCallbackActivityInstance struct {
	ID string `json:"id"`
}
