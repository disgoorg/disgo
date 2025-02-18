package slicehelper

import (
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func TestJoinSnowflakes(t *testing.T) {
	data := []struct {
		name       string
		snowflakes []snowflake.ID
		expected   string
	}{
		{
			name:       "0 snowflakes",
			snowflakes: []snowflake.ID{},
			expected:   "",
		},
		{
			name:       "1 snowflake",
			snowflakes: []snowflake.ID{snowflake.ID(1)},
			expected:   "1",
		},
		{
			name:       "2 snowflakes",
			snowflakes: []snowflake.ID{snowflake.ID(1), snowflake.ID(2)},
			expected:   "1,2",
		},
		{
			name:       "3 snowflakes",
			snowflakes: []snowflake.ID{snowflake.ID(1), snowflake.ID(2), snowflake.ID(3)},
			expected:   "1,2,3",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			actual := JoinSnowflakes(d.snowflakes)
			if actual != d.expected {
				t.Errorf("expected %s, got %s", d.expected, actual)
			}
		})
	}
}
