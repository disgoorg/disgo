package main

import (
	"errors"
	"os"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rpc"
)

var (
	channelID    = snowflake.GetEnv("disgo_channel_id")
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

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var tokenRs *discord.AccessTokenResponse
	code, err := client.Authorize([]discord.OAuth2Scope{discord.OAuth2ScopeRPC}, "", "")
	if err != nil {
		log.Fatal(err)
	}

	tokenRs, err = oauth2Client.GetAccessToken(clientID, clientSecret, code, "http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.Authenticate(tokenRs.AccessToken); err != nil {
		log.Fatal(err)
	}

	if channel, err := client.SelectVoiceChannel(channelID, false, false); err != nil {
		var dataError rpc.EventDataError
		if errors.As(err, &dataError) {
			if dataError.Code == 5003 { // User is in a voice channel, try again with force
				if channel, err = client.SelectVoiceChannel(channelID, true, false); err != nil {
					log.Fatal(err)
				} else {
					log.Info(channel)
				}
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Info(channel)
	}
}
