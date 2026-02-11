package discord

import (
	"fmt"
	"io"
	"slices"

	"github.com/disgoorg/snowflake/v2"
)

// NewWebhookMessageCreate returns a new WebhookMessageCreate with no fields set.
func NewWebhookMessageCreate() WebhookMessageCreate {
	return WebhookMessageCreate{}
}

// NewWebhookMessageCreateV2 returns a new WebhookMessageCreate with MessageFlagIsComponentsV2 flag set & allows to directly pass components.
func NewWebhookMessageCreateV2(components ...LayoutComponent) WebhookMessageCreate {
	return WebhookMessageCreate{
		Flags:      MessageFlagIsComponentsV2,
		Components: components,
	}
}

type WebhookMessageCreate struct {
	Content         string             `json:"content,omitempty"`
	Username        string             `json:"username,omitempty"`
	AvatarURL       string             `json:"avatar_url,omitempty"`
	TTS             bool               `json:"tts,omitempty"`
	Embeds          []Embed            `json:"embeds,omitempty"`
	Components      []LayoutComponent  `json:"components,omitempty"`
	Attachments     []AttachmentCreate `json:"attachments,omitempty"`
	Files           []*File            `json:"-"`
	AllowedMentions *AllowedMentions   `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags       `json:"flags,omitempty"`
	ThreadName      string             `json:"thread_name,omitempty"`
	AppliedTags     []snowflake.ID     `json:"applied_tags,omitempty"`
	Poll            *PollCreate        `json:"poll,omitempty"`
}

// ToBody returns the WebhookMessageCreate ready for body.
func (m WebhookMessageCreate) ToBody() (any, error) {
	if len(m.Files) > 0 {
		m.Attachments = parseAttachments(m.Files)
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

// WithContent returns a new WebhookMessageCreate with the provided content.
func (m WebhookMessageCreate) WithContent(content string) WebhookMessageCreate {
	m.Content = content
	return m
}

// WithContentf returns a new WebhookMessageCreate with the formatted content.
func (m WebhookMessageCreate) WithContentf(content string, a ...any) WebhookMessageCreate {
	m.Content = fmt.Sprintf(content, a...)
	return m
}

// WithUsername returns a new WebhookMessageCreate with the provided username.
func (m WebhookMessageCreate) WithUsername(username string) WebhookMessageCreate {
	m.Username = username
	return m
}

// WithAvatarURL returns a new WebhookMessageCreate with the provided avatar URL.
func (m WebhookMessageCreate) WithAvatarURL(url string) WebhookMessageCreate {
	m.AvatarURL = url
	return m
}

// WithTTS returns a new WebhookMessageCreate with the provided TTS setting.
func (m WebhookMessageCreate) WithTTS(tts bool) WebhookMessageCreate {
	m.TTS = tts
	return m
}

// WithEmbeds returns a new WebhookMessageCreate with the provided Embed(s).
func (m WebhookMessageCreate) WithEmbeds(embeds ...Embed) WebhookMessageCreate {
	m.Embeds = embeds
	return m
}

// WithEmbed returns a new WebhookMessageCreate with the provided Embed at the index.
func (m WebhookMessageCreate) WithEmbed(i int, embed Embed) WebhookMessageCreate {
	if len(m.Embeds) > i {
		m.Embeds = slices.Clone(m.Embeds)
		m.Embeds[i] = embed
	}
	return m
}

// AddEmbeds returns a new WebhookMessageCreate with the provided embeds added.
func (m WebhookMessageCreate) AddEmbeds(embeds ...Embed) WebhookMessageCreate {
	m.Embeds = append(m.Embeds, embeds...)
	return m
}

// ClearEmbeds returns a new WebhookMessageCreate with no embeds.
func (m WebhookMessageCreate) ClearEmbeds() WebhookMessageCreate {
	m.Embeds = []Embed{}
	return m
}

// RemoveEmbed returns a new WebhookMessageCreate with the embed at the index removed.
func (m WebhookMessageCreate) RemoveEmbed(i int) WebhookMessageCreate {
	if len(m.Embeds) > i {
		m.Embeds = slices.Delete(slices.Clone(m.Embeds), i, i+1)
	}
	return m
}

// WithComponents returns a new WebhookMessageCreate with the provided LayoutComponent(s).
func (m WebhookMessageCreate) WithComponents(components ...LayoutComponent) WebhookMessageCreate {
	m.Components = components
	return m
}

// UpdateComponent returns a new WebhookMessageCreate with the provided LayoutComponent at the index.
func (m WebhookMessageCreate) UpdateComponent(id int, component LayoutComponent) WebhookMessageCreate {
	for i, cc := range m.Components {
		if cc.GetID() == id {
			m.Components = slices.Clone(m.Components)
			m.Components[i] = component
			return m
		}
	}
	return m
}

// AddComponents returns a new WebhookMessageCreate with the provided LayoutComponent(s) added.
func (m WebhookMessageCreate) AddComponents(containers ...LayoutComponent) WebhookMessageCreate {
	m.Components = append(m.Components, containers...)
	return m
}

// AddActionRow returns a new WebhookMessageCreate with a new ActionRowComponent containing the provided InteractiveComponent(s) added.
func (m WebhookMessageCreate) AddActionRow(components ...InteractiveComponent) WebhookMessageCreate {
	m.Components = append(m.Components, NewActionRow(components...))
	return m
}

// RemoveComponent returns a new WebhookMessageCreate with the LayoutComponent at the index removed.
func (m WebhookMessageCreate) RemoveComponent(id int) WebhookMessageCreate {
	for i, cc := range m.Components {
		if cc.GetID() == id {
			m.Components = slices.Delete(slices.Clone(m.Components), i, i+1)
			return m
		}
	}
	return m
}

// ClearComponents returns a new WebhookMessageCreate with no LayoutComponent(s).
func (m WebhookMessageCreate) ClearComponents() WebhookMessageCreate {
	m.Components = []LayoutComponent{}
	return m
}

// WithFiles returns a new WebhookMessageCreate with the provided File(s).
func (m WebhookMessageCreate) WithFiles(files ...*File) WebhookMessageCreate {
	m.Files = files
	return m
}

// UpdateFile returns a new WebhookMessageCreate with the provided File at the index.
func (m WebhookMessageCreate) UpdateFile(i int, file *File) WebhookMessageCreate {
	if len(m.Files) > i {
		m.Files = slices.Clone(m.Files)
		m.Files[i] = file
	}
	return m
}

// AddFiles returns a new WebhookMessageCreate with the File(s) added.
func (m WebhookMessageCreate) AddFiles(files ...*File) WebhookMessageCreate {
	m.Files = append(m.Files, files...)
	return m
}

// AddFile returns a new WebhookMessageCreate with a File added.
func (m WebhookMessageCreate) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) WebhookMessageCreate {
	m.Files = append(m.Files, NewFile(name, description, reader, flags...))
	return m
}

// RemoveFile returns a new WebhookMessageCreate with the File at the index removed.
func (m WebhookMessageCreate) RemoveFile(i int) WebhookMessageCreate {
	if len(m.Files) > i {
		m.Files = slices.Delete(slices.Clone(m.Files), i, i+1)
	}
	return m
}

// ClearFiles returns a new WebhookMessageCreate with no File(s).
func (m WebhookMessageCreate) ClearFiles() WebhookMessageCreate {
	m.Files = []*File{}
	return m
}

// WithAllowedMentions returns a new WebhookMessageCreate with the provided AllowedMentions.
func (m WebhookMessageCreate) WithAllowedMentions(allowedMentions *AllowedMentions) WebhookMessageCreate {
	m.AllowedMentions = allowedMentions
	return m
}

// ClearAllowedMentions returns a new WebhookMessageCreate with no AllowedMentions.
func (m WebhookMessageCreate) ClearAllowedMentions() WebhookMessageCreate {
	return m.WithAllowedMentions(nil)
}

// WithFlags returns a new WebhookMessageCreate with the provided message flags.
func (m WebhookMessageCreate) WithFlags(flags ...MessageFlags) WebhookMessageCreate {
	m.Flags = m.Flags.Add(flags...)
	return m
}

// AddFlags returns a new WebhookMessageCreate with the provided MessageFlags added.
func (m WebhookMessageCreate) AddFlags(flags ...MessageFlags) WebhookMessageCreate {
	m.Flags = m.Flags.Add(flags...)
	return m
}

// RemoveFlags returns a new WebhookMessageCreate with the provided MessageFlags removed.
func (m WebhookMessageCreate) RemoveFlags(flags ...MessageFlags) WebhookMessageCreate {
	m.Flags = m.Flags.Remove(flags...)
	return m
}

// ClearFlags returns a new WebhookMessageCreate with no MessageFlags.
func (m WebhookMessageCreate) ClearFlags() WebhookMessageCreate {
	return m.WithFlags(MessageFlagsNone)
}

// WithIsComponentsV2 returns a new WebhookMessageCreate with MessageFlagIsComponentsV2 added/removed.
// Once a message with the flag has been sent, it cannot be removed by editing the message.
func (m WebhookMessageCreate) WithIsComponentsV2(isComponentV2 bool) WebhookMessageCreate {
	if isComponentV2 {
		m.Flags = m.Flags.Add(MessageFlagIsComponentsV2)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagIsComponentsV2)
	}
	return m
}

// WithSuppressEmbeds returns a new WebhookMessageCreate with MessageFlagSuppressEmbeds added/removed.
func (m WebhookMessageCreate) WithSuppressEmbeds(suppressEmbeds bool) WebhookMessageCreate {
	if suppressEmbeds {
		m.Flags = m.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	return m
}

// WithSuppressNotifications returns a new WebhookMessageCreate with MessageFlagSuppressNotifications added/removed.
func (m WebhookMessageCreate) WithSuppressNotifications(suppressNotifications bool) WebhookMessageCreate {
	if suppressNotifications {
		m.Flags = m.Flags.Add(MessageFlagSuppressNotifications)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagSuppressNotifications)
	}
	return m
}

// WithThreadName returns a new WebhookMessageCreate with the provided thread name.
func (m WebhookMessageCreate) WithThreadName(threadName string) WebhookMessageCreate {
	m.ThreadName = threadName
	return m
}

// WithAppliedTags returns a new WebhookMessageCreate with the provided applied tags.
func (m WebhookMessageCreate) WithAppliedTags(appliedTags ...snowflake.ID) WebhookMessageCreate {
	m.AppliedTags = appliedTags
	return m
}

// AddAppliedTags returns a new WebhookMessageCreate with the provided applied tags added.
func (m WebhookMessageCreate) AddAppliedTags(appliedTags ...snowflake.ID) WebhookMessageCreate {
	m.AppliedTags = append(m.AppliedTags, appliedTags...)
	return m
}

// RemoveAppliedTag returns a new WebhookMessageCreate with the provided applied tag removed.
func (m WebhookMessageCreate) RemoveAppliedTag(tagId snowflake.ID) WebhookMessageCreate {
	m.AppliedTags = slices.DeleteFunc(slices.Clone(m.AppliedTags), func(id snowflake.ID) bool {
		return id == tagId
	})
	return m
}

// ClearAppliedTags returns a new WebhookMessageCreate with no applied tags.
func (m WebhookMessageCreate) ClearAppliedTags() WebhookMessageCreate {
	m.AppliedTags = []snowflake.ID{}
	return m
}

// WithPoll returns a new WebhookMessageCreate with the provided Poll.
func (m WebhookMessageCreate) WithPoll(poll PollCreate) WebhookMessageCreate {
	m.Poll = &poll
	return m
}

// ClearPoll returns a new WebhookMessageCreate with no Poll.
func (m WebhookMessageCreate) ClearPoll() WebhookMessageCreate {
	m.Poll = nil
	return m
}
