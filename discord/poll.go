package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           time.Time      `json:"expiry"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type"`
	Results          []PollResults  `json:"results"`
}

type PollMedia struct {
	Text  *string       `json:"text"`
	Emoji *PartialEmoji `json:"emoji"`
}

type PollAnswer struct {
	AnswerID  *int      `json:"answer_id"`
	PollMedia PollMedia `json:"poll_media"`
}

type PollResults struct {
	IsFinalized  bool              `json:"is_finalized"`
	AnswerCounts []PollAnswerCount `json:"answer_counts"`
}

type PollAnswerCount struct {
	ID      snowflake.ID `json:"id"`
	Count   int          `json:"count"`
	MeVoted bool         `json:"me_voted"`
}

type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = iota + 1
)
