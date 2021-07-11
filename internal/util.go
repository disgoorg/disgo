package internal

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/DisgoOrg/disgo/api"
)

// IDFromToken returns the applicationID from the BotToken
func IDFromToken(token string) (*api.Snowflake, error) {
	strs := strings.Split(token, ".")
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

func normalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}
