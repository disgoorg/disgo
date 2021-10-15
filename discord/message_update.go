package discord

// MessageUpdate is used to edit a Message
type MessageUpdate struct {
	Content         *string          `json:"content,omitempty"`
	Embeds          *[]Embed         `json:"embeds,omitempty"`
	Components      *[]Component     `json:"components,omitempty"`
	Attachments     *[]Attachment    `json:"attachments,omitempty"`
	Files           []*File          `json:"-"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           *MessageFlags    `json:"flags,omitempty"`
}

func (_ MessageUpdate) dataType() dataType {
	return dataTypeMessageUpdate
}

// ToBody returns the MessageUpdate ready for body
func (m MessageUpdate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

func (m MessageUpdate) ToResponseBody(response InteractionResponse) (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(response, m.Files...)
	}
	return m, nil
}
