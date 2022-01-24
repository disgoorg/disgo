package discord

import "github.com/DisgoOrg/snowflake"

//Attachment is used for files sent in a Message
type Attachment struct {
	ID        snowflake.Snowflake `json:"id,omitempty"`
	Filename  string              `json:"filename"`
	Size      int                 `json:"size"`
	URL       string              `json:"url"`
	ProxyURL  string              `json:"proxy_url"`
	Height    *int                `json:"height"`
	Width     *int                `json:"width"`
	Ephemeral bool                `json:"ephemeral"`
}
