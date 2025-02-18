package slicehelper

import (
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

func JoinSnowflakes(snowflakes []snowflake.ID) string {
	strs := make([]string, len(snowflakes))
	for i, s := range snowflakes {
		strs[i] = s.String()
	}
	return strings.Join(strs, ",")
}
