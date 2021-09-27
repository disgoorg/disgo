package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

// ActivityType represents the status of a user, one of Game, Streaming, Listening, Watching, Custom or Competing
type ActivityType int

// Constants for Activity(s)
//goland:noinspection GoUnusedConst
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
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Type          ActivityType        `json:"type"`
	URL           *string             `json:"url"`
	CreatedAt     Time                `json:"created_at"`
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
	Buttons       []ActivityButton    `json:"buttons"`
}

func (a *Activity) UnmarshalJSON(data []byte) error {
	var v struct {
		Name          string              `json:"name"`
		Type          ActivityType        `json:"type"`
		URL           *string             `json:"url"`
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
		CreatedAt     int64               `json:"created_at"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Activity{
		Name:          v.Name,
		Type:          v.Type,
		URL:           v.URL,
		Timestamps:    v.Timestamps,
		ApplicationID: v.ApplicationID,
		Details:       v.Details,
		State:         v.State,
		Emoji:         v.Emoji,
		Party:         v.Party,
		Assets:        v.Assets,
		Secrets:       v.Secrets,
		Instance:      v.Instance,
		Flags:         v.Flags,
		CreatedAt:     Time{Time: time.Unix(v.CreatedAt/int64(time.Millisecond), (v.CreatedAt%int64(time.Millisecond))*int64(time.Nanosecond))},
	}
	return nil
}

func (a Activity) MarshalJSON() ([]byte, error) {
	v := struct {
		Name          string              `json:"name"`
		Type          ActivityType        `json:"type"`
		URL           *string             `json:"url"`
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
		CreatedAt     int64               `json:"created_at"`
	}{
		Name:          a.Name,
		Type:          a.Type,
		URL:           a.URL,
		Timestamps:    a.Timestamps,
		ApplicationID: a.ApplicationID,
		Details:       a.Details,
		State:         a.State,
		Emoji:         a.Emoji,
		Party:         a.Party,
		Assets:        a.Assets,
		Secrets:       a.Secrets,
		Instance:      a.Instance,
		Flags:         a.Flags,
		CreatedAt:     a.CreatedAt.UnixNano() / int64(time.Millisecond),
	}
	return json.Marshal(v)
}

type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ActivityTimestamps represents when a user started and ended their activity
type ActivityTimestamps struct {
	Start *Time `json:"start,omitempty"`
	End   *Time `json:"end,omitempty"`
}

func (t *ActivityTimestamps) UnmarshalJSON(data []byte) error {
	var v struct {
		Start *int64 `json:"start,omitempty"`
		End   *int64 `json:"end,omitempty"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Start != nil {
		t.Start = &Time{Time: time.Unix(*v.Start/int64(time.Millisecond), (*v.Start%int64(time.Millisecond))*int64(time.Nanosecond))}
	}
	if v.End != nil {
		t.End = &Time{Time: time.Unix(*v.End/int64(time.Millisecond), (*v.End%int64(time.Millisecond))*int64(time.Nanosecond))}
	}
	return nil
}

func (t ActivityTimestamps) MarshalJSON() ([]byte, error) {
	v := struct {
		Start *int64 `json:"start,omitempty"`
		End   *int64 `json:"end,omitempty"`
	}{}
	if t.Start != nil {
		start := t.Start.UnixNano() / int64(time.Millisecond)
		v.Start = &start
	}
	if t.End != nil {
		end := t.End.UnixNano() / int64(time.Millisecond)
		v.End = &end
	}
	return json.Marshal(v)
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
