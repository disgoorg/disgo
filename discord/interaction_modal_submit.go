package discord

import (
	"iter"

	"github.com/disgoorg/json/v2"
)

var (
	_ Interaction = (*ModalSubmitInteraction)(nil)
)

type ModalSubmitInteraction struct {
	baseInteraction
	Data ModalSubmitInteractionData `json:"data"`
	// Message is only present if the modal was triggered from a button
	Message *Message `json:"message,omitempty"`
}

func (i *ModalSubmitInteraction) UnmarshalJSON(data []byte) error {
	var interaction struct {
		rawInteraction
		Data    ModalSubmitInteractionData `json:"data"`
		Message *Message                   `json:"message,omitempty"`
	}
	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}

	i.baseInteraction.id = interaction.ID
	i.baseInteraction.applicationID = interaction.ApplicationID
	i.baseInteraction.token = interaction.Token
	i.baseInteraction.version = interaction.Version
	i.baseInteraction.guild = interaction.Guild
	i.baseInteraction.guildID = interaction.GuildID
	i.baseInteraction.channel = interaction.Channel
	i.baseInteraction.locale = interaction.Locale
	i.baseInteraction.guildLocale = interaction.GuildLocale
	i.baseInteraction.member = interaction.Member
	i.baseInteraction.user = interaction.User
	i.baseInteraction.appPermissions = interaction.AppPermissions
	i.baseInteraction.entitlements = interaction.Entitlements
	i.baseInteraction.authorizingIntegrationOwners = interaction.AuthorizingIntegrationOwners
	i.baseInteraction.context = interaction.Context
	i.baseInteraction.attachmentSizeLimit = interaction.AttachmentSizeLimit

	if i.baseInteraction.member != nil && i.baseInteraction.guildID != nil {
		i.baseInteraction.member.GuildID = *i.baseInteraction.guildID
	}

	i.Data = interaction.Data
	i.Message = interaction.Message
	if i.Message != nil {
		i.Message.GuildID = i.baseInteraction.guildID
	}
	return nil
}

func (i ModalSubmitInteraction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		rawInteraction
		Data ModalSubmitInteractionData `json:"data"`
		// Message is only present if the modal was triggered from a button
		Message *Message `json:"message,omitempty"`
	}{
		rawInteraction: rawInteraction{
			ID:                           i.id,
			Type:                         i.Type(),
			ApplicationID:                i.applicationID,
			Token:                        i.token,
			Version:                      i.version,
			Guild:                        i.guild,
			GuildID:                      i.guildID,
			Channel:                      i.channel,
			Locale:                       i.locale,
			GuildLocale:                  i.guildLocale,
			Member:                       i.member,
			User:                         i.user,
			AppPermissions:               i.appPermissions,
			Entitlements:                 i.entitlements,
			AuthorizingIntegrationOwners: i.authorizingIntegrationOwners,
			Context:                      i.context,
			AttachmentSizeLimit:          i.attachmentSizeLimit,
		},
		Data:    i.Data,
		Message: i.Message,
	})
}

func (ModalSubmitInteraction) Type() InteractionType {
	return InteractionTypeModalSubmit
}

func (ModalSubmitInteraction) interaction() {}

type ModalSubmitInteractionData struct {
	CustomID   string            `json:"custom_id"`
	Components []LayoutComponent `json:"components"`
}

func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	var iData struct {
		CustomID   string               `json:"custom_id"`
		Components []UnmarshalComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	d.CustomID = iData.CustomID

	components := make([]LayoutComponent, 0, len(iData.Components))
	for _, containerComponent := range iData.Components {
		components = append(components, containerComponent.Component.(LayoutComponent))
	}
	d.Components = components
	return nil
}

func (d ModalSubmitInteractionData) AllComponents() iter.Seq[Component] {
	return componentIter(d.Components)
}

func (d ModalSubmitInteractionData) Component(customID string) (InteractiveComponent, bool) {
	for component := range d.AllComponents() {
		if ic, ok := component.(InteractiveComponent); ok && ic.GetCustomID() == customID {
			return ic, true
		}
	}
	return nil, false
}

func (d ModalSubmitInteractionData) TextInputComponent(customID string) (TextInputComponent, bool) {
	if component, ok := d.Component(customID); ok {
		textInputComponent, ok := component.(TextInputComponent)
		return textInputComponent, ok
	}
	return TextInputComponent{}, false
}

func (d ModalSubmitInteractionData) OptText(customID string) (string, bool) {
	if textInputComponent, ok := d.TextInputComponent(customID); ok {
		return textInputComponent.Value, true
	}
	return "", false
}

func (d ModalSubmitInteractionData) Text(customID string) string {
	if text, ok := d.OptText(customID); ok {
		return text
	}
	return ""
}
