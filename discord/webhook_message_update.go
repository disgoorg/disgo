package discord

// WebhookMessageUpdate is used to edit a Message
type WebhookMessageUpdate struct {
	Content         *string             `json:"content,omitempty"`
	Embeds          *[]Embed            `json:"embeds,omitempty"`
	Components      *[]LayoutComponent  `json:"components,omitempty"`
	Attachments     *[]AttachmentUpdate `json:"attachments,omitempty"`
	Files           []*File             `json:"-"`
	AllowedMentions *AllowedMentions    `json:"allowed_mentions,omitempty"`
	Poll            *PollCreate         `json:"poll,omitempty"`
	// Flags are the MessageFlags of the message.
	// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
	// Use MessageFlags.Add for flags like discord.MessageFlagIsComponentsV2.
	Flags *MessageFlags `json:"flags,omitempty"`
}

// ToBody returns the WebhookMessageUpdate ready for body
func (m WebhookMessageUpdate) ToBody() (any, error) {
	if len(m.Files) > 0 {
		for _, attachmentCreate := range parseAttachments(m.Files) {
			if m.Attachments == nil {
				m.Attachments = new([]AttachmentUpdate)
			}
			*m.Attachments = append(*m.Attachments, attachmentCreate)
		}
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}
