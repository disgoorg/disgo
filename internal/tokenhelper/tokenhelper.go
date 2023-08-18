package tokenhelper

import (
	"encoding/base64"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// IDFromToken returns the applicationID from the token
func IDFromToken(token string) (*snowflake.ID, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return nil, discord.ErrInvalidBotToken
	}
	byteID, err := base64.RawStdEncoding.DecodeString(strs[0])
	if err != nil {
		return nil, err
	}
	strID, err := snowflake.Parse(string(byteID))
	if err != nil {
		return nil, err
	}
	return &strID, nil
}
