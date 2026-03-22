package slicehelper

import (
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

// JoinSnowflakes joins a slice of snowflake IDs into a comma-separated string.
func JoinSnowflakes(snowflakes []snowflake.ID) string {
	strs := make([]string, len(snowflakes))
	for i, s := range snowflakes {
		strs[i] = s.String()
	}
	return strings.Join(strs, ",")
}

// JoinStrings joins a slice of values with a ~string underlying type into a comma-separated string.
func JoinStrings[T ~string](vals []T) string {
	strs := make([]string, len(vals))
	for i, v := range vals {
		strs[i] = string(v)
	}
	return strings.Join(strs, ",")
}
