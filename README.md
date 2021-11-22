[![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/disgo.svg)](https://pkg.go.dev/github.com/DisgoOrg/disgo)
[![Go Report](https://goreportcard.com/badge/github.com/DisgoOrg/disgo)](https://goreportcard.com/report/github.com/DisgoOrg/disgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DisgoOrg/disgo)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgo/blob/master/LICENSE)
[![Disgo Version](https://img.shields.io/github/v/tag/DisgoOrg/disgo?label=release)](https://github.com/DisgoOrg/disgo/releases/latest)
[![Disgo Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/TewhTfDpvW)

<img align="right" src="/.github/discord_gopher.png" width=192 alt="discord gopher">

# disgo

disgo is a [Discord](https://discord.com) API wrapper written in [Go](https://golang.org/) aimed to be consistent, modular, customizable and easy to use

## Summary

1. [Features](#features)
2. [Missing Features](#missing-features)
3. [Getting Started](#getting-started)
4. [Documentation](#documentation)
5. [Examples](#examples)
6. [Related Projects](#related-projects)
7. [Why another library?](#why-another-library)
8. [Troubleshooting](#troubleshooting)
9. [Contributing](#contributing)
10. [License](#license)

### Features

* Full Rest API coverage
* [Gateway](https://discord.com/developers/docs/topics/gateway)
* [Sharding](https://discord.com/developers/docs/topics/gateway#sharding)
* [HTTP Interactions](https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction)
* [Application Commands](https://discord.com/developers/docs/interactions/application-commands)
* [Message Components](https://discord.com/developers/docs/interactions/message-components)
* [Stage Instance](https://discord.com/developers/docs/resources/stage-instance)
* [Guild Template](https://discord.com/developers/docs/resources/guild-template)
* [Sticker](https://discord.com/developers/docs/resources/sticker)
* [RateLimit](https://discord.com/developers/docs/topics/rate-limits)
* [Webhook](https://discord.com/developers/docs/resources/webhook)
* [OAuth2](https://discord.com/developers/docs/topics/oauth2)
* [Threads](https://discord.com/developers/docs/topics/threads)

### Missing Features

* [Voice](https://discord.com/developers/docs/topics/voice-connections)
* [RPC](https://discord.com/developers/docs/topics/rpc)
* [Guild Scheduled Event](https://discord.com/developers/docs/resources/guild-scheduled-event)

## Getting Started

### Installing

```sh
go get github.com/DisgoOrg/disgo
```

### Building a Disgo Instance

```go
package main

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

func main() {
	disgo, err := bot.New("token",
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(
				discord.GatewayIntentGuilds,
				discord.GatewayIntentGuildMessages,
				discord.GatewayIntentDirectMessages,
			),
		),
	)
}
```

### Full Ping Pong Example

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

func main() {
	disgo, err := bot.New(os.Getenv("token"),
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(
				discord.GatewayIntentGuilds,
				discord.GatewayIntentGuildMessages,
				discord.GatewayIntentDirectMessages,
			),
		),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	defer disgo.Close()

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	var message string
	if event.Message.Content == "ping" {
		message = "pong"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}
	if message != "" {
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}

```

### Logging

disgo uses our own small [logging lib](https://github.com/DisgoOrg/log) which provides an [interface](https://github.com/DisgoOrg/log/blob/master/logger.go) you can implement. This lib also comes with a default logger which is interchangeable and based on the standard log package. You can read more about it [here](https://github.com/DisgoOrg/log)

## Documentation

Documentation is wip and can be found under

* [![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/disgo.svg)](https://pkg.go.dev/github.com/DisgoOrg/disgo)
* [![Discord Documentation](https://img.shields.io/badge/Discord%20Documentation-blue.svg)](https://discord.com/developers/docs)

Wiki is currently under construction

## Examples

You can find examples under [_examples](https://github.com/DisgoOrg/disgo/tree/master/_examples)

or in these projects:

* [disgo-butler](https://github.com/DisgoOrg/disgo-butler)
* [BansBot](https://github.com/Skye-31/BansBot)
* [Reddit-Discord-Bot](https://github.com/TopiSenpai/Reddit-Discord-Bot)
* [Kitsune-Bot](https://github.com/TopiSenpai/Kitsune-Bot)
* [Uno-Bot](https://github.com/TopiSenpai/Uno-Bot)

## Related Projects

### [Lavalink](https://github.com/freyacodes/Lavalink)

Standalone audio sending node based on Lavaplayer and JDA-Audio. Allows for sending audio without it ever reaching any of your shards.

Being used in production by FredBoat, Dyno, LewdBot, and more.

### [disgolink](https://github.com/DisgoOrg/disgolink)

[Lavalink Client](https://github.com/freyacodes/Lavalink) which can be used to communicate with LavaLink to play/search tracks

### [disgofy](https://github.com/DisgoOrg/disgofy)

[disgolink](https://github.com/DisgoOrg/disgolink) Spotify integration. disgofy resolved Spotify urls to tracks/albums/playlists and lazy searches for them on YouTube

### [dislog](https://github.com/DisgoOrg/dislog)

Discord webhook logger integration for [logrus](https://github.com/sirupsen/logrus)

### [disgommand](https://github.com/DisgoOrg/disgommand)

Command framework for disgo in [gorilla/mux](https://github.com/gorilla/mux) style

## Why another library?

[discordgo](https://github.com/bwmarrin/discordgo) is a great library, but it's super low level and pain
[disgord](https://github.com/andersfylling/disgord) I don't like code gen magic
[arikawa](https://github.com/diamondburned/arikawa) v3 rewrite looks promising but when I started with disgo v2 looked kinda bad

disgo aims to be a high level library that is modular and not a pain to use.

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/TewhTfDpvW)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/TewhTfDpvW) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgo/blob/master/LICENSE)
. See LICENSE for more information.


