package main

import (
	_ "embed"
	"os"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/rest"
)

var token = os.Getenv("disgo_token")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Info("bot version: ", disgo.Version)

	client := rest.New(rest.NewClient(token))

	page := client.GetMessagesPage(817327182111571989, 1016790288607498240, 3)

	var i int
	for page.Next() {
		for _, m := range page.Items {
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
