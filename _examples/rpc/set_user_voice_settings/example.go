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
	channelID    = snowflake.GetEnv("disgo_channel_id")
	clientID     = snowflake.GetEnv("disgo_client_id")
	clientSecret = os.Getenv("disgo_client_secret")
	userID       = snowflake.GetEnv("disgo_user_id")
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

	var mute bool

	if channel, err := client.GetChannel(channelID); err != nil {
		log.Fatal(err)
	} else {
		if channel.Type != discord.ChannelTypeGuildVoice && channel.Type != discord.ChannelTypeGuildStageVoice {
			log.Fatal("Error: Selected channel is not voice channel")
		}
		var found = false

		for _, state := range channel.VoiceStates {
			if state.User.ID != userID {
				continue
			}
			found = true
			mute = !state.Mute
			break
		}
		if !found {
			log.Fatal("Error: Voice state for specified user not found.")
		}
	}

	settings := rpc.CmdArgsSetUserVoiceSettings{
		UserID: userID,
		Mute:   &mute,
	}

	if voiceSettings, err := client.SetUserVoiceSettings(settings); err != nil {
		log.Fatal(err)
	} else {
		log.Info(voiceSettings)
	}
}
