package discord

type WebhookMessageCreate struct {
	Content         string           `json:"content,omitempty"`
	Username        string           `json:"username,omitempty"`
	AvatarURL       string           `json:"avatar_url,omitempty"`
	TTS             bool             `json:"tts,omitempty"`
	Embeds          []Embed          `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	Files           []*File          `json:"-"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
}

// ToBody returns the MessageCreate ready for body
func (m WebhookMessageCreate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}
