package models

import (
	"time"

	"github.com/DiscoOrg/disgo/constants"
)

const ()

type Message struct {
	Reactions       []Reactions           `json:"reactions"`
	Attachments     []interface{}         `json:"attachments"`
	Tts             bool                  `json:"tts"`
	Embeds          []interface{}         `json:"embeds"`
	Timestamp       time.Time             `json:"timestamp"`
	MentionEveryone bool                  `json:"mention_everyone"`
	ID              Snowflake             `json:"id"`
	Pinned          bool                  `json:"pinned"`
	EditedTimestamp interface{}           `json:"edited_timestamp"`
	Author          Author                `json:"author"`
	MentionRoles    []interface{}         `json:"mention_roles"`
	Content         string                `json:"content"`
	ChannelID       Snowflake             `json:"channel_id"`
	Mentions        []interface{}         `json:"mentions"`
	ChannelType     constants.ChannelType `json:"type"`
}
type Emoji struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
}
type Reactions struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}
type Author struct {
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	ID            Snowflake `json:"id"`
	Avatar        string    `json:"avatar"`
}
