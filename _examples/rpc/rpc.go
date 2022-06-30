package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rpc"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var clientID = snowflake.GetEnv("disgo_client_id")

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Info("example is starting...")

	eventHandler := func(event rpc.Event, data rpc.MessageData) {
		//log.Infof("event: %s, data: %#v", event, data)
	}

	client, err := rpc.NewClient(clientID, eventHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	err = client.Send(rpc.CmdAuthorize, rpc.CmdArgsAuthorize{
		ClientID: clientID,
		Scopes:   []discord.OAuth2Scope{discord.OAuth2ScopeMessagesRead},
	}, rpc.CmdHandler(func(data rpc.CmdRsAuthorize) {
		println("handleAuthorize")
	}))

	if err != nil {
		log.Fatal(err)
	}

	//oauth2Client := rest.cl

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
