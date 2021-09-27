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
7. [Troubleshooting](#troubleshooting)
8. [Contributing](#contributing)
9. [License](#license)

### Features

* Full Rest API coverage
* [Gateway](https://discord.com/developers/docs/topics/gateway) support
* [Sharding](https://discord.com/developers/docs/topics/gateway#sharding) support
* [HTTP Interactions](https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction) support
* [Application Commands](https://discord.com/developers/docs/interactions/application-commands) support
* [Message Components](https://discord.com/developers/docs/interactions/message-components) support
* [Stage Instance](https://discord.com/developers/docs/resources/stage-instance) support
* [Guild Template](https://discord.com/developers/docs/resources/guild-template) support
* [Sticker](https://discord.com/developers/docs/resources/sticker) support
* [RateLimit](https://discord.com/developers/docs/topics/rate-limits) handling
* [Webhook](https://discord.com/developers/docs/resources/webhook) support
* [OAuth2](https://discord.com/developers/docs/topics/oauth2) support

### Missing Features

* [Voice](https://discord.com/developers/docs/resources/voice) support
* [Threads](https://discord.com/developers/docs/topics/threads) support

## Getting Started

### Installing

```sh
go get github.com/DisgoOrg/disgo
```

### Building a Disgo Instance

```go
disgo, err := core.NewBot(os.Getenv("token"),
    // set which gateway intents we should use
    core.WithGatewayConfigOpts(
        gateway.WithGatewayIntents(
            discord.GatewayIntentGuilds,
            discord.GatewayIntentGuildMessages,
            discord.GatewayIntentDirectMessages,
        ),
    ),
	// set what to cache
    core.WithCacheConfigOpts(
		core.WithCacheFlags(core.CacheFlagsDefault),
    ),
    // add our event listeners
    core.WithEventListeners(&core.ListenerAdapter{
        OnMessageCreate: onMessageCreate,
    }),
)
if err != nil {
    log.Fatal("error while building disgo: ", err)
}

// connect to the gateway
if err := disgo.ConnectGateway(); err != nil {
    log.Fatal("error while connecting to the gateway: ", err)
}
```

### Full Ping Pong Example

```go
package main

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/DisgoOrg/disgo/core"
    "github.com/DisgoOrg/disgo/discord"
    "github.com/DisgoOrg/disgo/gateway"
    "github.com/DisgoOrg/log"
)

func main() {
    disgo, err := core.NewBot(os.Getenv("token"),
        core.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
        core.WithCacheConfigOpts(core.WithCacheFlags(core.CacheFlagGuilds)),
        core.WithEventListeners(&core.ListenerAdapter{
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

func onMessageCreate(event *core.MessageCreateEvent) {
    _, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
```

### Logging

disgo uses our own small [logging lib](https://github.com/DisgoOrg/log) which provides an [interface](https://github.com/DisgoOrg/log/blob/master/logger.go) you can implement. This lib also comes with a default logger which is interchangeable and based on the standard log package. You can read more about it [here](https://github.com/DisgoOrg/log)

## Documentation

Documentation is unfinished and can be found under

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

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/TewhTfDpvW)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/TewhTfDpvW) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgo/blob/master/LICENSE)
. See LICENSE for more information.


