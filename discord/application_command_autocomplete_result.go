package discord

type AutocompleteResult struct {
	Choices []AutocompleteChoice `json:"choices"`
}

func (_ AutocompleteResult) interactionCallbackData() {}

type AutocompleteChoice interface {
	autoCompleteChoice()
}

type AutocompleteChoiceString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (_ AutocompleteChoiceString) autoCompleteChoice() {}

type AutocompleteChoiceInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (_ AutocompleteChoiceInt) autoCompleteChoice() {}

type AutocompleteChoiceFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (_ AutocompleteChoiceFloat) autoCompleteChoice() {}
