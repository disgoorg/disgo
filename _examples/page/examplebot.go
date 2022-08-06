package main

import (
	_ "embed"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/log"
)

var token = os.Getenv("disgo_token")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Info("bot version: ", disgo.Version)

	client := rest.New(rest.NewClient(token))

	page := client.GetMessagesPage(817327182111571989, 1003431540228882432, 0, 3)

	var i int
	for page.Previous() {
		for _, m := range page.Data {
			println(m.ID)
		}
		println("---")
		i++
		if i >= 3 {
			break
		}
	}
	if page.Err != nil {
		log.Error(page.Err)
	}
}
