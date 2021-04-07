package api

import "time"

// ActivityType represents the status of a user, one of Game, Streaming, Listening, Watching, Custom or Competing
type ActivityType int

// Constants for activities
const (
	Game ActivityType = iota
	Streaming
	Listening
	Watching
	Custom
	Competing
)

// Activity represents the fields of a user's presence
type Activity struct {
	Name          string              `json:"name"`
	Type          ActivityType        `json:"type"`
	URL           *string             `json:"url"`
	CreatedAt     time.Time           `json:"created_at"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID Snowflake           `json:"application_id,omitempty"`
	Details       *string             `json:"details,omitempty"`
	State         *string             `json:"state,omitempty"`
	Emoji         *ActivityEmoji      `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      *bool               `json:"instance,omitempty"`
	Flags         int                 `json:"flags,omitempty"`
}

// ActivityTimestamps represents when a user started and ended their activity
type ActivityTimestamps struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// ActivityEmoji is an Emoji object for an Activity
type ActivityEmoji struct {
	Name     string     `json:"name"`
	ID       *Snowflake `json:"id,omitempty"`
	Animated *bool      `json:"animated,omitempty"`
}

// ActivityParty is information about the party of the player
type ActivityParty struct {
	ID   Snowflake `json:"id,omitempty"`
	Size []int     `json:"size,omitempty"`
}

// ActivityAssets are the images for the presence and hover texts
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// ActivitySecrets contain secrets for Rich Presence joining and spectating
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}
