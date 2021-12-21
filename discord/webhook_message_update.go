package discord

// WebhookMessageUpdate is used to edit a Message
type WebhookMessageUpdate struct {
	Content         *string               `json:"content,omitempty"`
	Embeds          *[]Embed              `json:"embeds,omitempty"`
	Components      *[]ContainerComponent `json:"components,omitempty"`
	Attachments     *[]Attachment         `json:"attachments,omitempty"`
	Files           []*File               `json:"-"`
	AllowedMentions *AllowedMentions      `json:"allowed_mentions,omitempty"`
}

// ToBody returns the WebhookMessageUpdate ready for body
func (m WebhookMessageUpdate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}
