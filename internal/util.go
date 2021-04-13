package internal

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
)

// IDFromToken returns the applicationID from the BotToken
func IDFromToken(token endpoints.Token) (*api.Snowflake, error) {
	strs := strings.Split(string(token), ".")
	if len(strs) == 0 {
		return nil, errors.New("BotToken is not in a valid format")
	}
	byteID, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return nil, err
	}
	strID := api.Snowflake(byteID)
	return &strID, nil
}
