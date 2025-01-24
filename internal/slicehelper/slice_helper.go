package slicehelper

import "github.com/disgoorg/snowflake/v2"

func JoinSnowflakes(snowflakes []snowflake.ID) string {
	var str string
	for i, s := range snowflakes {
		str += s.String()
		if i != len(str)-1 {
			str += ","
		}
	}
	return str
}
