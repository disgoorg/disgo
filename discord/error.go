package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
)

type ErrorCode int

var _ error = (*APIError)(nil)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Errors  string    `json:"errors"`
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var v struct {
		Code    ErrorCode       `json:"code"`
		Message string          `json:"message"`
		Errors  json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("error unmarshalling discord error: %w", err)
	}

	e.Code = v.Code
	e.Message = v.Message
	e.Errors = string(v.Errors)
	return nil
}

func (e *APIError) Error() string {
	return e.Message
}
