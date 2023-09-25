package main

import (
	"log/slog"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/rest"
)

var token = os.Getenv("disgo_token")

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client := rest.New(rest.NewClient(token))

	page := client.GetMessagesPage(817327182111571989, 1016790288607498240, 3)

	var i int
	for page.Next() {
		for _, m := range page.Items {
			slog.Info(m.ID.String())
		}
		slog.Info("---")
		i++
		if i >= 3 {
			break
		}
	}
	if page.Err != nil {
		slog.Error("error getting messages", slog.Any("err", page.Err))
	}
}
