package discord

import (
	"slices"
	"time"

	"github.com/disgoorg/json/v2"
)

// NewPollCreate returns a new PollCreate with the given question and answers.
func NewPollCreate(question string, answers ...PollMedia) PollCreate {
	return PollCreate{
		Question: PollMedia{
			Text: &question,
		},
		Answers: answers,
	}
}

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

// WithQuestion sets the question of the Poll
func (p PollCreate) WithQuestion(text string) PollCreate {
	p.Question = PollMedia{
		Text: &text,
	}
	return p
}

// SetAnswers sets the answers of the Poll
func (p PollCreate) WithAnswers(answers ...PollMedia) PollCreate {
	p.Answers = answers
	return p
}

// AddAnswer adds an answer to the Poll
func (p PollCreate) AddAnswer(text string, emoji *PartialEmoji) PollCreate {
	p.Answers = append(p.Answers, PollMedia{
		Text:  &text,
		Emoji: emoji,
	})
	return p
}

// RemoveAnswer removes an answer from the Poll
func (p PollCreate) RemoveAnswer(i int) PollCreate {
	if len(p.Answers) > i {
		p.Answers = slices.Delete(slices.Clone(p.Answers), i, i+1)
	}
	return p
}

// ClearAnswers removes all answers of the Poll
func (p PollCreate) ClearAnswers() PollCreate {
	p.Answers = []PollMedia{}
	return p
}

// WithDuration sets the duration of the Poll (in hours)
func (p PollCreate) WithDuration(duration int) PollCreate {
	p.Duration = duration
	return p
}

// WithAllowMultiselect sets whether users will be able to vote for more than one answer of the Poll
func (p PollCreate) WithAllowMultiselect(multiselect bool) PollCreate {
	p.AllowMultiselect = multiselect
	return p
}

// WithLayoutType sets the layout of the Poll
func (p PollCreate) WithLayoutType(layout PollLayoutType) PollCreate {
	p.LayoutType = layout
	return p
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
