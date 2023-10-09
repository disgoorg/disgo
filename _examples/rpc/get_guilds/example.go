package main

import (
	"fmt"
	"os"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rpc"
)

var (
	clientID     = snowflake.GetEnv("disgo_client_id")
	clientSecret = os.Getenv("disgo_client_secret")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Info("example is starting...")

	oauth2Client := rest.NewOAuth2(rest.NewClient(""))

	client, err := rpc.NewClient(clientID)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	var tokenRs *discord.AccessTokenResponse
	code, err := client.Authorize(rpc.CmdArgsAuthorize{
		ClientID: clientID,
		Scopes:   []discord.OAuth2Scope{discord.OAuth2ScopeRPC, discord.OAuth2ScopeGuilds},
	})
	if err != nil {
		log.Fatal(err)
	}

	tokenRs, err = oauth2Client.GetAccessToken(clientID, clientSecret, code, "http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.Authenticate(rpc.CmdArgsAuthenticate{AccessToken: tokenRs.AccessToken}); err != nil {
		log.Fatal(err)
	}

	if guilds, err := client.GetGuilds(); err != nil {
		log.Fatal(err)
	} else {
		for _, guild := range guilds {
			if guild.IconURL == nil {
				log.Info(guild.Name)
			} else {
				log.Info(fmt.Sprintf("%s: %s", guild.Name, *guild.IconURL))
			}
		}
	}
}
