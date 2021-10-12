package discord

type AutocompleteResult struct {
	Choices []AutocompleteChoice `json:"choices"`
}

func (_ AutocompleteResult) dataType() dataType {
	return dataTypeAutocompleteResult
}

func (m AutocompleteResult) ToBody() (interface{}, error) {
	return m, nil
}

type AutocompleteChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
