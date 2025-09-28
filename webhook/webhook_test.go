package webhook

import (
	"testing"

	"github.com/disgoorg/snowflake/v2"
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
			c, err := NewWithURL(tc.URL)
			if tc.Err {
				if err == nil {
					t.Errorf("expected error for URL %q, got none", tc.URL)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for URL %q: %v", tc.URL, err)
			}
			if c.ID != tc.ID {
				t.Errorf("expected ID %v, got %v", tc.ID, c.ID)
			}
			if c.Token != tc.Token {
				t.Errorf("expected Token %q, got %q", tc.Token, c.Token)
			}
		})
	}
}
