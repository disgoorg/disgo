package discord

import (
	"time"

	"github.com/disgoorg/json"
)

type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           *time.Time     `json:"expiry"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type"`
	Results          *PollResults   `json:"results"`
}

type PollCreate struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollMedia    `json:"-"`
	Duration         int            `json:"duration"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
}

func (p PollCreate) MarshalJSON() ([]byte, error) {
	type pollCreate PollCreate

	answers := make([]PollAnswer, 0, len(p.Answers))
	for _, answer := range p.Answers {
		answers = append(answers, PollAnswer{
			PollMedia: answer,
		})
	}
	return json.Marshal(struct {
		Answers []PollAnswer `json:"answers"`
		pollCreate
	}{
		Answers:    answers,
		pollCreate: pollCreate(p),
	})
}

type PollMedia struct {
	Text  *string       `json:"text"`
	Emoji *PartialEmoji `json:"emoji,omitempty"`
}

type PollAnswer struct {
	AnswerID  *int      `json:"answer_id,omitempty"`
	PollMedia PollMedia `json:"poll_media"`
}

type PollResults struct {
	IsFinalized  bool              `json:"is_finalized"`
	AnswerCounts []PollAnswerCount `json:"answer_counts"`
}

type PollAnswerCount struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = iota + 1
)
