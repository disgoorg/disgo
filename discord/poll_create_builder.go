package discord

// PollCreateBuilder helps create PollCreate structs easier
type PollCreateBuilder struct {
	PollCreate
}

// SetQuestion sets the question of the Poll
func (b *PollCreateBuilder) SetQuestion(text string) *PollCreateBuilder {
	b.Question = PollMedia{
		Text: &text,
	}
	return b
}

// SetAnswers sets the answers of the Poll
func (b *PollCreateBuilder) SetAnswers(answers ...PollMedia) *PollCreateBuilder {
	b.Answers = answers
	return b
}

// AddAnswer adds an answer to the Poll
func (b *PollCreateBuilder) AddAnswer(text string, emoji *PartialEmoji) *PollCreateBuilder {
	b.Answers = append(b.Answers, PollMedia{
		Text:  &text,
		Emoji: emoji,
	})
	return b
}

// RemoveAnswer removes an answer from the Poll
func (b *PollCreateBuilder) RemoveAnswer(i int) *PollCreateBuilder {
	if len(b.Answers) > i {
		b.Answers = append(b.Answers[:i], b.Answers[i+1:]...)
	}
	return b
}

// ClearAnswers removes all answers of the Poll
func (b *PollCreateBuilder) ClearAnswers() *PollCreateBuilder {
	b.Answers = []PollMedia{}
	return b
}

// SetDuration sets the duration of the Poll (in hours)
func (b *PollCreateBuilder) SetDuration(duration int) *PollCreateBuilder {
	b.Duration = duration
	return b
}

// SetAllowMultiselect sets whether users will be able to vote for more than one answer of the Poll
func (b *PollCreateBuilder) SetAllowMultiselect(multiselect bool) *PollCreateBuilder {
	b.AllowMultiselect = multiselect
	return b
}

// SetLayoutType sets the layout of the Poll
func (b *PollCreateBuilder) SetLayoutType(layout PollLayoutType) *PollCreateBuilder {
	b.LayoutType = layout
	return b
}

// Build builds the PollCreateBuilder to a PollCreate struct
func (b *PollCreateBuilder) Build() PollCreate {
	return b.PollCreate
}
