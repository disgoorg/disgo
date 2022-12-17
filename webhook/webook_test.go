package webhook

import (
	"testing"

	"github.com/disgoorg/snowflake/v2"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tt := []struct {
		URL   string
		ID    snowflake.ID
		Token string
		Err   bool
	}{
		{
			URL:   "https://discord.com/api/webhooks/123456789123456789/foo",
			ID:    snowflake.ID(123456789123456789),
			Token: "foo",
		},
		{
			URL:   "https://discord.com/api/webhooks/123456789123456789/foo/",
			ID:    snowflake.ID(123456789123456789),
			Token: "foo",
		},
		{
			URL:   "https://canary.discord.com/api/webhooks/123456789123456789/foo",
			ID:    snowflake.ID(123456789123456789),
			Token: "foo",
		},
		{
			URL: "foobarbaz",
			Err: true,
		},
		{
			URL:   "https://discord.com/api/webhooks/123456789123456789/foo?wait=10",
			ID:    snowflake.ID(123456789123456789),
			Token: "foo",
		},
	}

	for _, tc := range tt {
		t.Run(tc.URL, func(t *testing.T) {
			assert := require.New(t)

			c, err := NewWithURL(tc.URL)
			if tc.Err {
				assert.Error(err, "URL parsing should have resulted in an error")
				return
			}
			assert.Equal(tc.ID, c.ID(), "URL ID should match")
			assert.Equal(tc.Token, c.Token(), "URL token should match")
		})
	}
}
