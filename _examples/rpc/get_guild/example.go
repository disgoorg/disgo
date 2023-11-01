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
	guildId      = snowflake.GetEnv("disgo_guild_id")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Info("example is starting...")

	oauth2Client := rest.NewOAuth2(rest.NewClient(""))

	client, err := rpc.New(clientID)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var tokenRs *discord.AccessTokenResponse
	codeRs, err := client.Authorize([]discord.OAuth2Scope{discord.OAuth2ScopeRPC, discord.OAuth2ScopeGuilds}, "", "")
	if err != nil {
		log.Fatal(err)
	}

	tokenRs, err = oauth2Client.GetAccessToken(clientID, clientSecret, codeRs.Code, "http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.Authenticate(tokenRs.AccessToken); err != nil {
		log.Fatal(err)
	}

	guild, err := client.GetGuild(guildId, 0)
	if err != nil {
		log.Fatal(err)
	}
	if guild.IconURL == nil {
		log.Info(guild.Name)
		return
	}
	log.Info(fmt.Sprintf("%s: %s", guild.Name, *guild.IconURL))
}
