package core

import (
	"encoding/base64"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
)

// IDFromToken returns the applicationID from the BotToken
//goland:noinspection GoUnusedExportedFunction
func IDFromToken(token string) (*discord.Snowflake, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return nil, discord.ErrInvalidBotToken
	}
	byteID, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return nil, err
	}
	strID := discord.Snowflake(byteID)
	return &strID, nil
}

