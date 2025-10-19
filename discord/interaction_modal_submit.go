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
	Resolved   ResolvedData      `json:"resolved"`
}

func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	var iData struct {
		CustomID   string               `json:"custom_id"`
		Components []UnmarshalComponent `json:"components"`
		Resolved   ResolvedData         `json:"resolved"`
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

	d.Resolved = iData.Resolved
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

func (d ModalSubmitInteractionData) TextInput(customID string) (TextInputComponent, bool) {
	if component, ok := d.Component(customID); ok {
		textInputComponent, ok := component.(TextInputComponent)
		return textInputComponent, ok
	}
	return TextInputComponent{}, false
}

func (d ModalSubmitInteractionData) OptText(customID string) (string, bool) {
	if textInputComponent, ok := d.TextInput(customID); ok {
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

func (d ModalSubmitInteractionData) StringSelectMenu(customID string) (StringSelectMenuComponent, bool) {
	if component, ok := d.Component(customID); ok {
		selectMenuComponent, ok := component.(StringSelectMenuComponent)
		return selectMenuComponent, ok
	}
	return StringSelectMenuComponent{}, false
}

func (d ModalSubmitInteractionData) OptStringValues(customID string) ([]string, bool) {
	if selectMenuComponent, ok := d.StringSelectMenu(customID); ok {
		return selectMenuComponent.Values, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) StringValues(customID string) []string {
	if text, ok := d.OptStringValues(customID); ok {
		return text
	}
	return nil
}

func (d ModalSubmitInteractionData) UserSelectMenu(customID string) (UserSelectMenuComponent, bool) {
	if component, ok := d.Component(customID); ok {
		selectMenuComponent, ok := component.(UserSelectMenuComponent)
		return selectMenuComponent, ok
	}
	return UserSelectMenuComponent{}, false
}

func (d ModalSubmitInteractionData) OptUsers(customID string) ([]User, bool) {
	if selectMenuComponent, ok := d.UserSelectMenu(customID); ok {
		users := make([]User, 0, len(selectMenuComponent.Values))
		for _, userID := range selectMenuComponent.Values {
			if user, ok := d.Resolved.Users[userID]; ok {
				users = append(users, user)
			}
		}
		return users, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) Users(customID string) []User {
	if users, ok := d.OptUsers(customID); ok {
		return users
	}
	return nil
}

func (d ModalSubmitInteractionData) RoleSelectMenu(customID string) (RoleSelectMenuComponent, bool) {
	if component, ok := d.Component(customID); ok {
		selectMenuComponent, ok := component.(RoleSelectMenuComponent)
		return selectMenuComponent, ok
	}
	return RoleSelectMenuComponent{}, false
}

func (d ModalSubmitInteractionData) OptRoles(customID string) ([]Role, bool) {
	if selectMenuComponent, ok := d.RoleSelectMenu(customID); ok {
		roles := make([]Role, 0, len(selectMenuComponent.Values))
		for _, roleID := range selectMenuComponent.Values {
			if role, ok := d.Resolved.Roles[roleID]; ok {
				roles = append(roles, role)
			}
		}
		return roles, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) Roles(customID string) []Role {
	if roles, ok := d.OptRoles(customID); ok {
		return roles
	}
	return nil
}

// MentionableValue is an interface for all values a [MentionableSelectMenuComponent] can return.
// [User]
// [ResolvedMember]
// [Role]
// [ResolvedChannel]
type MentionableValue interface {
	isMentionableValue()
}

func (d ModalSubmitInteractionData) MentionableSelectMenu(customID string) (MentionableSelectMenuComponent, bool) {
	if component, ok := d.Component(customID); ok {
		selectMenuComponent, ok := component.(MentionableSelectMenuComponent)
		return selectMenuComponent, ok
	}
	return MentionableSelectMenuComponent{}, false
}

func (d ModalSubmitInteractionData) OptMentionables(customID string) ([]Mentionable, bool) {
	if selectMenuComponent, ok := d.MentionableSelectMenu(customID); ok {
		mentionables := make([]Mentionable, 0, len(selectMenuComponent.Values))
		for _, id := range selectMenuComponent.Values {
			if user, ok := d.Resolved.Users[id]; ok {
				mentionables = append(mentionables, user)
				continue
			}
			if role, ok := d.Resolved.Roles[id]; ok {
				mentionables = append(mentionables, role)
				continue
			}
		}
		return mentionables, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) Mentionables(customID string) []Mentionable {
	if mentionables, ok := d.OptMentionables(customID); ok {
		return mentionables
	}
	return nil
}

func (d ModalSubmitInteractionData) ChannelSelectMenu(customID string) (ChannelSelectMenuComponent, bool) {
	if component, ok := d.Component(customID); ok {
		selectMenuComponent, ok := component.(ChannelSelectMenuComponent)
		return selectMenuComponent, ok
	}
	return ChannelSelectMenuComponent{}, false
}

func (d ModalSubmitInteractionData) OptChannels(customID string) ([]ResolvedChannel, bool) {
	if selectMenuComponent, ok := d.ChannelSelectMenu(customID); ok {
		channels := make([]ResolvedChannel, 0, len(selectMenuComponent.Values))
		for _, channelID := range selectMenuComponent.Values {
			if channel, ok := d.Resolved.Channels[channelID]; ok {
				channels = append(channels, channel)
			}
		}
		return channels, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) Channels(customID string) []ResolvedChannel {
	if channels, ok := d.OptChannels(customID); ok {
		return channels
	}
	return nil
}

func (d ModalSubmitInteractionData) FileUpload(customID string) (FileUploadComponent, bool) {
	if component, ok := d.Component(customID); ok {
		fileUploadComponent, ok := component.(FileUploadComponent)
		return fileUploadComponent, ok
	}
	return FileUploadComponent{}, false
}

func (d ModalSubmitInteractionData) OptAttachments(customID string) ([]Attachment, bool) {
	if fileUploadComponent, ok := d.FileUpload(customID); ok {
		attachments := make([]Attachment, 0, len(fileUploadComponent.Values))
		for _, attachmentID := range fileUploadComponent.Values {
			if attachment, ok := d.Resolved.Attachments[attachmentID]; ok {
				attachments = append(attachments, attachment)
			}
		}
		return attachments, true
	}
	return nil, false
}

func (d ModalSubmitInteractionData) Attachments(customID string) []Attachment {
	if attachments, ok := d.OptAttachments(customID); ok {
		return attachments
	}
	return nil
}
