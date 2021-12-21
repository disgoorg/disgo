package discord

type AutocompleteResult struct {
	Choices []AutocompleteChoice `json:"choices"`
}

func (AutocompleteResult) interactionCallbackData() {}

type AutocompleteChoice interface {
	autoCompleteChoice()
}

type AutocompleteChoiceString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (AutocompleteChoiceString) autoCompleteChoice() {}

type AutocompleteChoiceInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (AutocompleteChoiceInt) autoCompleteChoice() {}

type AutocompleteChoiceFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (AutocompleteChoiceFloat) autoCompleteChoice() {}
