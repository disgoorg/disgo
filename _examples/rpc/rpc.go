package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rpc"
)

var (
	clientID     = snowflake.GetEnv("disgo_client_id")
	clientSecret = os.Getenv("disgo_client_secret")
	channelID    = snowflake.GetEnv("disgo_channel_id")
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
	if err = client.Send(rpc.Message{
		Cmd: rpc.CmdAuthorize,
		Args: rpc.CmdArgsAuthorize{
			ClientID: clientID,
			Scopes:   []discord.OAuth2Scope{discord.OAuth2ScopeRPC, discord.OAuth2ScopeMessagesRead},
		},
	}, rpc.NewHandler(func(data rpc.CmdRsAuthorize) {
		tokenRs, err = oauth2Client.GetAccessToken(clientID, clientSecret, data.Code, "http://localhost")
		if err != nil {
			log.Fatal(err)
		}
	})); err != nil {
		log.Fatal(err)
	}

	if err = client.Send(rpc.Message{
		Cmd: rpc.CmdAuthenticate,
		Args: rpc.CmdArgsAuthenticate{
			AccessToken: tokenRs.AccessToken,
		},
	}, nil); err != nil {
		log.Fatal(err)
	}

	if err = client.Subscribe(rpc.EventMessageCreate, rpc.CmdArgsSubscribeMessage{
		ChannelID: channelID,
	}, rpc.NewHandler(func(data rpc.EventDataMessageCreate) {
		log.Info("message: ", data.Message.Content)
	})); err != nil {
		log.Fatal(err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
