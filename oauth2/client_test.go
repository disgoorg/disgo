package oauth2

import (
	"net/url"
	"strings"
	"testing"

	"github.com/disgoorg/disgo/discord"
)

func TestClient_GenerateAuthorizationURL_Permissions(t *testing.T) {
	client := New(1234567890, "secret")

	authURL := client.GenerateAuthorizationURL(AuthorizationURLParams{
		RedirectURI: "https://example.com/callback",
		Scopes:      []discord.OAuth2Scope{discord.OAuth2ScopeBot},
		Permissions: discord.PermissionSendMessages | discord.PermissionStream,
	})

	_, query, _ := strings.Cut(authURL, "?")
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("unexpected error parsing query: %v", err)
	}

	expected := "2560"
	if got := values.Get("permissions"); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
