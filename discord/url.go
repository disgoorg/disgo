package discord

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

type QueryValues map[string]any

func (q QueryValues) Encode() string {
	values := url.Values{}
	for k, v := range q {
		values.Set(k, fmt.Sprint(v))
	}
	return values.Encode()
}

func urlPrint(url string, params ...any) string {
	for _, param := range params {
		start := strings.Index(url, "{")
		end := strings.Index(url, "}")
		url = url[:start] + fmt.Sprint(param) + url[end+1:]
	}
	return url
}

func InviteURL(code string) string {
	return urlPrint("https://discord.gg/{code}", code)
}

func WebhookURL(webhookID snowflake.ID, webhookToken string) string {
	return urlPrint("https://discord.com/api/webhooks/{webhook.id}/{webhook.token}", webhookID, webhookToken)
}

func AuthorizeURL(values QueryValues) string {
	query := values.Encode()
	if query != "" {
		query = "?" + query
	}
	return "https://discord.com/api/oauth2/authorize" + query
}
