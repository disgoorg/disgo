package main

import (
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

	var deaf bool

	if voiceSettings, err := client.GetVoiceSettings(); err != nil {
		log.Fatal(err)
	} else {
		deaf = !voiceSettings.Deaf
	}

	settings := rpc.CmdArgsSetVoiceSettings{
		Deaf: &deaf,
	}

	if voiceSettings, err := client.SetVoiceSettings(settings); err != nil {
		log.Fatal(err)
	} else {
		log.Info(voiceSettings)
	}
}