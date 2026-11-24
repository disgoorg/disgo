package rest

import (
	"testing"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

func TestThreadMemberPage_Next(t *testing.T) {
	var requestedAfter snowflake.ID
	p := &ThreadMemberPage{
		getItems: func(after snowflake.ID) ([]discord.ThreadMember, error) {
			requestedAfter = after
			return []discord.ThreadMember{
				{UserID: 100},
				{UserID: 200},
				{UserID: 300},
			}, nil
		},
	}

	if !p.Next() {
		t.Fatal("expected Next to succeed")
	}

	// Calling Next() again should use the LAST item's ID (300) as cursor, not the first (100)
	if !p.Next() {
		t.Fatal("expected second Next to succeed")
	}
	if requestedAfter != 300 {
		t.Fatalf("expected after ID to be 300, got %d", requestedAfter)
	}
}

func TestPollAnswerVotesPage_Next(t *testing.T) {
	var requestedAfter snowflake.ID
	p := &PollAnswerVotesPage{
		getItems: func(after snowflake.ID) ([]discord.User, error) {
			requestedAfter = after
			return []discord.User{
				{ID: 10},
				{ID: 20},
				{ID: 30},
			}, nil
		},
	}

	if !p.Next() {
		t.Fatal("expected Next to succeed")
	}

	if !p.Next() {
		t.Fatal("expected second Next to succeed")
	}
	if requestedAfter != 30 {
		t.Fatalf("expected after ID to be 30, got %d", requestedAfter)
	}
}
