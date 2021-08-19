package util

import (
	"encoding/base64"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
)

// IDFromToken returns the applicationID from the BotToken
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

func NormalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}

// GetOS returns the simplified version of the operating system for sending to Discord in the IdentifyCommandDataProperties. OS payload
func GetOS() string {
	OS := runtime.GOOS
	if strings.HasPrefix(OS, "windows") {
		return "windows"
	}
	if strings.HasPrefix(OS, "darwin") {
		return "darwin"
	}
	return "linux"
}

// GitHub is the Disgo GitHub URL
const GitHub = "https://github.com/DisgoOrg/disgo"

// Version returns the current used Disgo version in the format vx.x.x
var Version = getVersion()

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(GitHub, dep.Path) {
				return dep.Version
			}
		}
	}
	return "unknown"
}
