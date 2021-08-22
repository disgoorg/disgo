package rest

import "strings"

//goland:noinspection GoUnusedExportedFunction
func NormalizeEmoji(emoji string) string {
	return strings.Replace(emoji, "#", "%23", -1)
}

