package rest

import "github.com/disgoorg/json"

// ItemsResponse serves as a utility to flatten objects wrapped in an "items" object into a slice of objects.
type ItemsResponse[T any] []T

func (rs *ItemsResponse[T]) UnmarshalJSON(data []byte) error {
	var response struct {
		Items []T `json:"items"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	*rs = response.Items
	return nil
}
